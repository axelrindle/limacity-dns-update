package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	IPV4 = "https://ifconfig.me"
	IPV6 = "https://ifconfig.co"
)

func MakeRequest(method string, base string, body io.Reader, args ...any) *http.Request {
	apiUrl := Env("API_URL", "https://www.lima-city.de/usercp")

	url := fmt.Sprintf(apiUrl+"/"+base, args...)
	req, _ := http.NewRequest(method, url, body)
	authString := Env("API_USER", "") + ":" + Env("API_PASSWORD", "")

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(authString)))

	return req
}

func LoadIPAddress(client *http.Client, ipType string) string {
	resp, err := client.Get(ipType)
	DieIf(err)

	defer resp.Body.Close()

	address, err := io.ReadAll(resp.Body)
	DieIf(err)

	return strings.Replace(string(address), "\n", "", 1)
}

func LoadDNSEntries(client *http.Client) []NameserverRecord {
	domainId := Env("DOMAIN_ID", "")
	request := MakeRequest("GET", "domains/%s/records.json", nil, domainId)

	resp, err := client.Do(request)
	DieIf(err)

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	DieIf(err)

	var records ResponseListRecords
	json.Unmarshal(body, &records)

	log.Printf("Loaded a total of %d DNS entries.\n", len(records.Records))

	return records.Records
}

func UpdateDNSRecord(client *http.Client, record NameserverRecord, addressType string, addressUrl string) {
	ipAddress := LoadIPAddress(client, addressUrl)

	record.Content = ipAddress
	log.Printf("Updating DNS entry '%s' with %s address '%s'\n",
		record.Name, addressType, record.Content)

	requestBody := RequestUpdateRecord{
		Record: record,
	}

	jsonString, err := json.Marshal(requestBody)
	DieIf(err)

	bodyReader := bytes.NewReader(jsonString)

	domainId := Env("DOMAIN_ID", "")
	request := MakeRequest("PUT", "domains/%s/records/%d", bodyReader, domainId, record.ID)
	request.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(request)
	DieIf(err)

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	DieIf(err)

	// var status ResponseUpdateRecord
	// json.Unmarshal(body, &status)

	log.Println(string(body))
}

func UpdateDNSv4Record(client *http.Client, record NameserverRecord) {
	UpdateDNSRecord(client, record, "IPV4", IPV4)
}

func UpdateDNSv6Record(client *http.Client, record NameserverRecord) {
	UpdateDNSRecord(client, record, "IPV6", IPV6)
}
