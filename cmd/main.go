package main

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/victor.diditskiy/replication_experiment/pkg/handlers/alive"
)

func main() {
	log := logrus.New()
	healthCheckHandler := alive.New(log)

	router := mux.NewRouter()
	router.Path("/alive").Handler(healthCheckHandler)

	log.Info("Starting web server")
	err := http.ListenAndServe(":80", router)
	if err != nil {
		log.Fatal(err)
	}
}
