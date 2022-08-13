package main

import (
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/victor_diditskiy/replication_experiment/pkg/data_generator"
	"github.com/victor_diditskiy/replication_experiment/pkg/dbpool"
	"github.com/victor_diditskiy/replication_experiment/pkg/handlers/alive"
	"github.com/victor_diditskiy/replication_experiment/pkg/handlers/generator"
	"github.com/victor_diditskiy/replication_experiment/pkg/storage"
)

const (
	DBConfigPath = "./config/db.yml"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	pool, err := dbpool.NewPool(DBConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	db := pool.GetLeader()
	_, err = db.Exec("insert into data (name, value) values ('test', 12345)")
	if err != nil {
		log.Fatal(err)
	}

	dbStorage := storage.New(pool)
	dataGenerator := data_generator.New(log, dbStorage)

	healthCheckHandler := alive.New(log)
	generatorHandler := generator.New(log, dataGenerator)

	router := mux.NewRouter()
	router.Path("/alive").Handler(healthCheckHandler)
	router.Path("/api/generate_data").Handler(generatorHandler)

	log.Info("Starting web server")
	err = http.ListenAndServe(":80", router)
	if err != nil {
		log.Fatal(err)
	}
}
