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
	BatchSize   int64 `json:"batch_size"`
	Limit       int64 `json:"limit"`
	ScaleFactor int64 `json:"scale_factor"`
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

	if req.BatchSize == 0 {
		message := "no batch size set"
		h.log.Error(message)
		w.WriteHeader(400)
		_, _ = w.Write([]byte("{\"status\": \"error\", \"message\": \"" + message + "\"}"))
		return
	}

	if req.Limit == 0 {
		message := "no limit set"
		h.log.Error(message)
		w.WriteHeader(400)
		_, _ = w.Write([]byte("{\"status\": \"error\", \"message\": \"" + message + "\"}"))
		return
	}

	if req.ScaleFactor == 0 {
		message := "no scale factor set"
		h.log.Error(message)
		w.WriteHeader(400)
		_, _ = w.Write([]byte("{\"status\": \"error\", \"message\": \"" + message + "\"}"))
		return
	}

	err = h.generator.Generate(r.Context(), req.BatchSize, req.Limit, req.ScaleFactor)
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
