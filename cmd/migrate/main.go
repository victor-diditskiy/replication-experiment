package main

import (
	"flag"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/victor_diditskiy/replication_experiment/pkg/dbpool"
	internalMigrator "github.com/victor_diditskiy/replication_experiment/pkg/migrator"
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

	migrator := internalMigrator.New(pool)

	var direction = flag.String("direction", "", "migration direction")
	var steps = flag.Int("steps", 0, "migration steps")
	flag.Parse()

	switch *direction {
	case "up":
		err = migrator.Up()
		if err != nil {
			log.Fatal(errors.Wrap(err, "failed to up migrations"))
		}
	case "down":
		err = migrator.Down() // or m.Step(2) if you want to explicitly set the number of migrations to run
		if err != nil {
			log.Fatal(errors.Wrap(err, "failed to down migrations"))
		}
	case "steps":
		err = migrator.Steps(*steps) // or m.Step(2) if you want to explicitly set the number of migrations to run
		if err != nil {
			log.Fatal(errors.Wrap(err, "failed to apply migration steps"))
		}

	default:
		log.Fatal("invalid direction flag given")
	}
}
