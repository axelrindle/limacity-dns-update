package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	client := &http.Client{Timeout: 10 * time.Second}
	gocron.Every(1).Hour().Do(task, client)

	log.Println("Task will run once every hour.")
	println()

	gocron.RunAll()
	<-gocron.Start()
}

func task(client *http.Client) {
	records := LoadDNSEntries(client)

	v4 := strings.Split(Env("DNS_IDS_IPV4", ","), ",")
	v6 := strings.Split(Env("DNS_IDS_IPV6", ""), ",")

	for _, record := range records {
		if SliceContains(v4, strconv.Itoa(record.ID)) {
			UpdateDNSv4Record(client, record)
		} else if SliceContains(v6, strconv.Itoa(record.ID)) {
			UpdateDNSv6Record(client, record)
		}
	}

	_, nextRun := gocron.NextRun()
	log.Printf("Next run at %s", nextRun.Local().Format(time.RFC3339))
	println()
}
