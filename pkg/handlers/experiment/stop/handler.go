package stop

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/victor_diditskiy/replication_experiment/pkg/plan"
)

type Handler struct {
	planManager *plan.Manager
	log         logrus.FieldLogger
}

func New(log logrus.FieldLogger, manager *plan.Manager) *Handler {
	return &Handler{
		planManager: manager,
		log:         log,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	command := plan.Command{
		Name: plan.StopCommand,
	}
	err := h.planManager.Execute(command)
	if err != nil {
		h.log.Error(fmt.Sprintf("failed to execute stop command: %s", err))
		w.WriteHeader(422)
		_, _ = w.Write([]byte("{\"status\": \"error\", \"message\": \"failed to execute stop command\"}"))
		return
	}

	h.log.WithFields(logrus.Fields{
		"command": plan.StartCommand,
	}).Info("plan successfully stopped")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, _ = w.Write([]byte("{\"status\": \"ok\"}"))
}
