package start

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/victor_diditskiy/replication_experiment/pkg/plan"
)

type Handler struct {
	planManager *plan.Manager
	log         logrus.FieldLogger
}

type ExperimentRequest struct {
	PlanName     string `json:"plan_name"`
	ReadWorkload *struct {
		ScaleFactor int `json:"scale_factor"`
	} `json:"read_workload,omitempty"`
	InsertWorkload *struct {
		ScaleFactor int `json:"scale_factor"`
	} `json:"insert_workload,omitempty"`
	UpdateWorkload *struct {
		ScaleFactor int `json:"scale_factor"`
	} `json:"update_workload,omitempty"`
}

func New(log logrus.FieldLogger, manager *plan.Manager) *Handler {
	return &Handler{
		planManager: manager,
		log:         log,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	req := &ExperimentRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		h.log.Error(fmt.Sprintf("failed to parse request: %s", err))
		w.WriteHeader(400)
		_, _ = w.Write([]byte("{\"status\": \"error\", \"message\": \"failed to parse request\"}"))
		return
	}

	config := plan.Config{}

	var planName plan.Name
	switch req.PlanName {
	case "read_only":
		planName = plan.ReadOnlyPlanName
	case "write_only":
		planName = plan.WriteOnlyPlanName
	case "combined":
		planName = plan.CombinedPlanName
	}
	if planName == "" {
		message := fmt.Sprintf("invalid request: invalid plan name %s", req.PlanName)
		h.log.Error(message)
		w.WriteHeader(400)
		_, _ = w.Write([]byte("{\"status\": \"error\", \"message\": \"" + message + "\"}"))
		return
	}

	if req.ReadWorkload != nil {
		config.ReadWorkload = &plan.ConfigItem{ScaleFactor: req.ReadWorkload.ScaleFactor}
	}

	if req.InsertWorkload != nil {
		config.InsertWorkload = &plan.ConfigItem{ScaleFactor: req.InsertWorkload.ScaleFactor}
	}

	if req.UpdateWorkload != nil {
		config.UpdateWorkload = &plan.ConfigItem{ScaleFactor: req.UpdateWorkload.ScaleFactor}
	}

	command := plan.Command{
		Name:     plan.StartCommand,
		PlanName: planName,
		Config:   config,
	}
	err = h.planManager.Execute(command)
	if err != nil {
		h.log.Error(fmt.Sprintf("failed to execute command: %s", err))
		w.WriteHeader(422)
		_, _ = w.Write([]byte("{\"status\": \"error\", \"message\": \"failed to execute command\"}"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, _ = w.Write([]byte("{\"status\": \"ok\"}"))
}
