package cmd

import (
	"context"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/axelrindle/limacity-dns-update/mock"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	godotenv.Load(".env.testing")
	logrus.SetOutput(io.Discard)

	server := mock.StartMock()

	exitCode := m.Run()

	server.Shutdown(context.Background())

	os.Exit(exitCode)
}

var httpClient = &http.Client{Timeout: 10 * time.Second}

func TestLoadDNSEntries(t *testing.T) {
	records, err := loadDNSEntries(httpClient)

	assert.Nil(t, err)
	assert.Len(t, records, 2, "Expected to receive 2 DNS records!")
}

func TestUpdate(t *testing.T) {
	records, err1 := loadDNSEntries(httpClient)

	assert.Nil(t, err1)
	assert.Len(t, records, 2, "Expected to receive 2 DNS records!")

	theRecord := records[0]

	assert.Equal(t, theRecord.ID, 123456)
	assert.Equal(t, theRecord.Content, "10.11.12.13")

	theRecord.Content = "13.12.11.10"

	err2 := handleRecord(httpClient, theRecord)

	assert.Nil(t, err2)

	records, err3 := loadDNSEntries(httpClient)

	assert.Nil(t, err3)
	assert.Len(t, records, 2, "Expected to receive 2 DNS records!")

	theRecordNew := records[0]

	assert.Equal(t, theRecordNew.ID, 123456)
	assert.Equal(t, theRecordNew.Content, "13.12.11.10")
}
