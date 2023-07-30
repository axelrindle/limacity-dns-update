package mock

import (
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func StartMock() *http.Server {
	server := createHTTPServer()

	go func() {
		log.WithField("address", server.Addr).Info("Server listening.")

		err := server.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			log.Info("Server closed.")
		} else if err != nil {
			log.WithField("error", err).Fatal("Server startup failed!")
		}
	}()

	return server
}
