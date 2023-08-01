package cmd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/axelrindle/limacity-dns-update/shared"
	log "github.com/sirupsen/logrus"
)

func makeRequest(method string, base string, body io.Reader, args ...any) *http.Request {
	apiUrl := shared.Env("API_URL", "https://www.lima-city.de/usercp")

	url := fmt.Sprintf(apiUrl+"/"+base, args...)
	req, _ := http.NewRequest(method, url, body)
	authString := shared.Env("API_USER", "") + ":" + shared.Env("API_PASSWORD", "")

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(authString)))

	return req
}

func loadIPAddress(client *http.Client, ipType string) (string, error) {
	resp, err := client.Get(ipType)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	address, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.Replace(string(address), "\n", "", 1), nil
}

func loadDNSEntries(client *http.Client) ([]shared.NameserverRecord, error) {
	domainId := shared.Env("DOMAIN_ID", "")
	request := makeRequest("GET", "domains/%s/records.json", nil, domainId)

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var records shared.ResponseListRecords
	errJson := json.Unmarshal(body, &records)
	if errJson != nil {
		return nil, errJson
	}

	log.WithFields(log.Fields{
		"count": len(records.Records),
	}).Info("Received DNS entries.")

	return records.Records, nil
}

var addressUrlDefaults = map[string]string{
	"API_URL_IPv4": "https://ifconfig.me/ip",
	"API_URL_IPv6": "https://ifconfig.co/ip",
}

func updateDNSRecord(client *http.Client, record shared.NameserverRecord, addressType string) error {
	addressUrlKey := "API_URL_" + addressType
	addressUrl := shared.Env(addressUrlKey, addressUrlDefaults[addressUrlKey])
	ipAddress, err := loadIPAddress(client, addressUrl)
	if err != nil {
		return err
	}

	logger := log.WithFields(log.Fields{
		"type":   addressType,
		"domain": record.Name,
		"record": record.Type,
	})

	if record.Content == ipAddress {
		logger.Info("Record already up-to-date.")
		return nil
	}

	record.Content = ipAddress
	requestBody := shared.RequestUpdateRecord{
		Record: record,
	}

	jsonString, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	bodyReader := bytes.NewReader(jsonString)

	domainId := shared.Env("DOMAIN_ID", "")
	request := makeRequest("PUT", "domains/%s/records/%d", bodyReader, domainId, record.ID)
	request.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var status shared.ResponseUpdateRecord
	json.Unmarshal(body, &status)

	if resp.StatusCode == 200 {
		logger.WithFields(log.Fields{
			"ip": ipAddress,
		}).Info("Update succeeded.")
	} else {
		log.WithFields(log.Fields{
			"error": status.Error,
		}).Error("Update failed!")
	}

	return nil
}

func updateDNSv4Record(client *http.Client, record shared.NameserverRecord) error {
	return updateDNSRecord(client, record, "IPv4")
}

func updateDNSv6Record(client *http.Client, record shared.NameserverRecord) error {
	return updateDNSRecord(client, record, "IPv6")
}
