package alive

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type HealthCheckHandler struct {
	log logrus.FieldLogger
}

func New(log logrus.FieldLogger) *HealthCheckHandler {
	return &HealthCheckHandler{
		log: log,
	}
}

func (h HealthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Alive")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, _ = w.Write([]byte("{\"status\": \"alive\"}"))
}
