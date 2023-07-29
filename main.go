package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	client := &http.Client{Timeout: 10 * time.Second}
	scheduler := gocron.NewScheduler(time.Local)
	scheduler.Every(1).Hour().Do(func() {
		task(client)

		_, nextRun := scheduler.NextRun()
		log.Printf("Next run at %s", nextRun.Local().Format(time.RFC3339))
		println()
	})

	log.Println("Task will run once every hour.")
	println()

	scheduler.StartBlocking()
}

func HandleRecord(client *http.Client, record NameserverRecord) error {
	v4 := strings.Split(Env("DNS_IDS_IPV4", ","), ",")
	v6 := strings.Split(Env("DNS_IDS_IPV6", ","), ",")

	if SliceContains(v4, strconv.Itoa(record.ID)) {
		return UpdateDNSv4Record(client, record)
	} else if SliceContains(v6, strconv.Itoa(record.ID)) {
		return UpdateDNSv6Record(client, record)
	}

	return nil
}

const fileFailure = "/tmp/failure"

func task(client *http.Client) {
	records, err := LoadDNSEntries(client)
	if err != nil {
		log.Println(err)
	} else {
		isContainer := Env("CONTAINER", "false")
		for _, record := range records {
			err := HandleRecord(client, record)
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
