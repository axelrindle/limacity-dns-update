package cmd

import (
	"context"
	"flag"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/axelrindle/limacity-dns-update/mock"
	"github.com/axelrindle/limacity-dns-update/shared"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"github.com/lnquy/cron"
	log "github.com/sirupsen/logrus"
)

const VERSION = "v0.3.0"

var invocation = 0

var showVersion bool
var oneshot bool
var startMock bool

func handleFlags() {
	flag.BoolVar(&showVersion, "version", false, "show the binary version and exit")
	flag.BoolVar(&oneshot, "oneshot", false, "run the updater once and exit")
	flag.BoolVar(&startMock, "mock", false, "Start the mock server")
	flag.Parse()

	if showVersion {
		println(VERSION)
		os.Exit(0)
	}

	if startMock {
		server := mock.StartMock()

		<-shared.GracefulShutdown(context.Background(), 5*time.Second, map[string]shared.ShutdownHook{
			"server": func(ctx context.Context) error {
				server.Shutdown(ctx)
				return nil
			},
		})
		os.Exit(0)
	}
}

func GetVersion() string {
	return VERSION
}

func Run() {
	handleFlags()

	godotenv.Load()

	if shared.Env("LOGGING_JSON", "false") == "true" {
		log.SetFormatter(&log.JSONFormatter{})
	}
	if shared.Env("LOGGING_DEBUG", "false") == "true" {
		log.SetLevel(log.DebugLevel)
	}

	client := &http.Client{Timeout: 10 * time.Second}

	if oneshot {
		task(client)
		return
	}

	scheduler := gocron.NewScheduler(time.Local)
	expression := shared.Env("CRON", "0 1 * * *")
	digits := len(strings.Split(expression, " "))
	handler := func() {
		invocation++
		log.WithFields(log.Fields{
			"invocation": invocation,
		}).Info("Starting task...")

		task(client)

		_, nextRun := scheduler.NextRun()
		log.WithFields(log.Fields{
			"timestamp": nextRun.Local().Format(time.RFC3339),
		}).Info("Done. Next run scheduled.")
		println()
	}

	switch digits {
	case 5:
		scheduler.Cron(expression).Do(handler)
	case 6:
		scheduler.CronWithSeconds(expression).Do(handler)
	default:
		panic("Invalid cron expression!")
	}

	locale := cron.Locale_en
	descriptor, err := cron.NewDescriptor(
		cron.Use24HourTimeFormat(true),
		cron.DayOfWeekStartsAtOne(true),
		cron.SetLocales(locale),
	)
	if err != nil {
		panic(err)
	}
	description, err := descriptor.ToDescription(expression, locale)
	if err != nil {
		panic(err)
	}

	runInitial := shared.Env("INITIAL", "true") == "true"
	if runInitial {
		log.Info("Initial execution")
		handler()
	}

	log.Info("Task schedule: " + description)
	println()

	scheduler.StartAsync()

	<-shared.GracefulShutdown(context.Background(), 2*time.Second, map[string]shared.ShutdownHook{
		"scheduler": func(ctx context.Context) error {
			scheduler.Stop()
			return nil
		},
	})
}

func handleRecord(client *http.Client, record shared.NameserverRecord) error {
	v4 := strings.Split(shared.Env("DNS_IDS_IPV4", ","), ",")
	v6 := strings.Split(shared.Env("DNS_IDS_IPV6", ","), ",")

	if shared.SliceContains(v4, strconv.Itoa(record.ID)) {
		return updateDNSv4Record(client, record)
	} else if shared.SliceContains(v6, strconv.Itoa(record.ID)) {
		return updateDNSv6Record(client, record)
	}

	return nil
}

const fileFailure = "/tmp/failure"

func task(client *http.Client) {
	records, err := loadDNSEntries(client)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to load DNS entries!")
	} else {
		isContainer := shared.Env("CONTAINER", "false")
		for _, record := range records {
			err := handleRecord(client, record)
			if isContainer == "false" {
				continue
			}

			if err != nil {
				os.WriteFile(fileFailure, nil, 0644)
			} else {
				os.Remove(fileFailure)
			}
		}
	}
}
