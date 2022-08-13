package generator

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/victor_diditskiy/replication_experiment/pkg/data_generator"
)

type Handler struct {
	generator *data_generator.Generator
	log       logrus.FieldLogger
}

type GenerateRequest struct {
	Limit int64 `json:"limit"`
}

func New(log logrus.FieldLogger, generator *data_generator.Generator) *Handler {
	return &Handler{
		generator: generator,
		log:       log,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	req := &GenerateRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		h.log.Error(fmt.Sprintf("failed to parse request: %s", err))
		w.WriteHeader(400)
		_, _ = w.Write([]byte("{\"status\": \"error\", \"message\": \"failed to parse request\"}"))
		return
	}

	err = h.generator.Generate(r.Context(), req.Limit)
	if err != nil {
		h.log.Error(fmt.Sprintf("failed to generate data: %s", err))
		w.WriteHeader(422)
		_, _ = w.Write([]byte("{\"status\": \"error\", \"message\": \"failed to generate data\"}"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, _ = w.Write([]byte("{\"status\": \"ok\"}"))
}
