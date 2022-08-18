package migrator

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/pkg/errors"

	"github.com/victor_diditskiy/replication_experiment/pkg/dbpool"
)

type Migrator struct {
	dbPool *dbpool.DBPool
}

func New(dbPool *dbpool.DBPool) *Migrator {
	return &Migrator{dbPool: dbPool}
}

func (m *Migrator) Up() error {
	for _, db := range m.dbPool.GetLeaders() {
		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			return errors.Wrap(err, "failed to init postgres wrapper")
		}

		migrator, err := migrate.NewWithDatabaseInstance(
			"file://migrations",
			"postgres", driver,
		)

		err = migrator.Up() // or m.Step(2) if you want to explicitly set the number of migrations to run
		if err != nil {
			return errors.Wrap(err, "failed to up migrations")
		}
	}

	return nil
}

func (m *Migrator) Down() error {
	for _, db := range m.dbPool.GetLeaders() {
		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			return errors.Wrap(err, "failed to init postgres wrapper")
		}

		migrator, err := migrate.NewWithDatabaseInstance(
			"file://migrations",
			"postgres", driver,
		)

		err = migrator.Down() // or m.Step(2) if you want to explicitly set the number of migrations to run
		if err != nil {
			return errors.Wrap(err, "failed to up migrations")
		}
	}

	return nil
}

func (m *Migrator) Steps(num int) error {
	for _, db := range m.dbPool.GetLeaders() {
		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			return errors.Wrap(err, "failed to init postgres wrapper")
		}

		migrator, err := migrate.NewWithDatabaseInstance(
			"file://migrations",
			"postgres", driver,
		)

		err = migrator.Steps(num)
		if err != nil {
			return errors.Wrapf(err, "failed to apply %d steps to migrations", num)
		}
	}

	return nil
}
