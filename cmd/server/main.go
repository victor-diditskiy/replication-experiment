package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"github.com/victor_diditskiy/replication_experiment/pkg/data_generator"
	"github.com/victor_diditskiy/replication_experiment/pkg/dbpool"
	"github.com/victor_diditskiy/replication_experiment/pkg/handlers/alive"
	"github.com/victor_diditskiy/replication_experiment/pkg/handlers/experiment/start"
	"github.com/victor_diditskiy/replication_experiment/pkg/handlers/experiment/stop"
	"github.com/victor_diditskiy/replication_experiment/pkg/handlers/generator"
	"github.com/victor_diditskiy/replication_experiment/pkg/plan"
	"github.com/victor_diditskiy/replication_experiment/pkg/storage"
	"github.com/victor_diditskiy/replication_experiment/pkg/workload"
)

const (
	DBConfigPath = "./config/db.yml"
	defaultPort  = 80
)

func main() {
	log := logrus.New()
	log.SetReportCaller(true)
	log.SetFormatter(&logrus.JSONFormatter{})

	pool, err := dbpool.NewPool(DBConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	dbStorage := storage.New(pool)
	metricStorage := storage.NewMetricStorage(dbStorage)
	dataGenerator := data_generator.New(log, metricStorage)

	workloads := workload.NewWorkloads(log, metricStorage)
	planManager := plan.NewManager(workloads)

	healthCheckHandler := alive.New(log)
	generatorHandler := generator.New(log, dataGenerator)
	startExperimentHandler := start.New(log, planManager)
	stopExperimentHandler := stop.New(log, planManager)

	router := mux.NewRouter()
	router.Path("/alive").Handler(healthCheckHandler).Methods("GET")
	router.Path("/api/generate_data").Handler(generatorHandler).Methods("GET")
	router.Path("/api/experiment/start").Handler(startExperimentHandler).Methods("POST")
	router.Path("/api/experiment/stop").Handler(stopExperimentHandler).Methods("POST")
	router.Path("/prometheus").Handler(promhttp.Handler())

	d := http.DefaultServeMux
	router.PathPrefix("/debug").Handler(d)

	log.Info("Starting web server")

	port := defaultPort
	portEnv := os.Getenv("PORT")
	if portEnv != "" {
		parsedPort, err := strconv.ParseInt(portEnv, 10, 64)
		if err == nil {
			port = int(parsedPort)
		}
	}
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		log.Fatal(err)
	}
}
