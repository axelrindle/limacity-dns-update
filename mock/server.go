package mock

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/axelrindle/limacity-dns-update/shared"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var records = []shared.NameserverRecord{
	{
		ID:      123456,
		Name:    "cloud.example.org",
		Type:    "A",
		Content: "10.11.12.13",
		TTL:     1800,
	},
	{
		ID:      654321,
		Name:    "cloud.example.org",
		Type:    "AAAA",
		Content: "8ace:c049:64ba:a47d:4a4f:87ba:116d:9efe",
	},
}

func addHandlers(router *mux.Router) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := []Route{}

		router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			path, err := route.GetPathTemplate()

			if err != nil {
				return err
			} else {
				response = append(response, Route{
					Path: path,
				})
			}

			return nil
		})

		data, err := json.Marshal(response)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Write(data)
		}
	}).Methods("GET")

	router.HandleFunc("/domains/{domainId}/records.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := shared.ResponseListRecords{
			Records: records,
		}

		data, err := json.Marshal(response)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Write(data)
		}
	}).Methods("GET")

	router.HandleFunc("/domains/{domainId}/records/{recordId}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)
		recordId := vars["recordId"]

		found := false
		for _, record := range records {
			if strconv.Itoa(record.ID) == recordId {
				found = true
				data, err := json.Marshal(record)
				if err != nil {
					w.Write([]byte(err.Error()))
					w.WriteHeader(http.StatusInternalServerError)
				} else {
					w.Write(data)
				}
			}
		}

		if !found {
			w.WriteHeader(http.StatusNotFound)
		}
	}).Methods("GET")

	router.HandleFunc("/domains/{domainId}/records/{recordId}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		recordId := vars["recordId"]

		found := false
		for index, record := range records {
			if strconv.Itoa(record.ID) == recordId {
				found = true

				body, err := io.ReadAll(r.Body)
				if err != nil {
					w.Write([]byte(err.Error()))
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				var newRecord RequestUpdateRecord
				json.Unmarshal(body, &newRecord)

				records[index] = newRecord.Record
			}
		}

		if !found {
			w.WriteHeader(http.StatusNotFound)
		}
	}).Methods("PUT")

	router.HandleFunc("/ip/v4", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("13.12.11.10"))
	})

	router.HandleFunc("/ip/v6", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("13cc:b6e4:3ec6:b0b0:1871:cf0d:af4c:3f6b"))
	})
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		log.WithFields(log.Fields{
			"method":    r.Method,
			"protocol":  r.Proto,
			"remote":    r.RemoteAddr,
			"type":      "request",
			"timestamp": startTime.Local().Format(time.RFC3339),
		}).Info(r.RequestURI)

		customWriter := &LogResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(customWriter, r)

		log.WithFields(log.Fields{
			"duration": time.Since(startTime).String(),
			"status":   customWriter.statusCode,
			"type":     "response",
		}).Info(r.RequestURI)
	})
}

func createHTTPServer() *http.Server {
	mux := mux.NewRouter()

	mux.Use(requestLogger)

	addHandlers(mux)

	return &http.Server{
		Addr:    ":3000",
		Handler: mux,
	}
}
