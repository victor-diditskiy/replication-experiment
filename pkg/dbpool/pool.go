package dbpool

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Leaders []struct {
		DSN string
	}
}

type DBPool struct {
	leaders   []*sql.DB
	followers []*sql.DB
}

func NewPool(dbConfigPath string) (*DBPool, error) {
	file, err := os.ReadFile(dbConfigPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read db config file")
	}

	conf := &Config{}
	err = yaml.Unmarshal(file, conf)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse db config file")
	}

	pool := &DBPool{}
	for i, leader := range conf.Leaders {
		if leader.DSN == "" {
			return nil, fmt.Errorf("no DSN set for leader #%d", i)
		}

		db, err := createDB(leader.DSN)

		if err != nil {
			return nil, errors.Wrap(err, "failed to init DB connection")
		}

		pool.leaders = append(pool.leaders, db)
	}

	return pool, nil
}

func createDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "failed to ping db")
	}

	return db, nil
}

func (p *DBPool) GetLeader() *sql.DB {
	// TODO: configure returning random leader
	return p.leaders[0]
}

func (p *DBPool) AllInstances() []*sql.DB {
	inst := make([]*sql.DB, 0, len(p.leaders)+len(p.followers))
	inst = append(inst, p.leaders...)
	inst = append(inst, p.followers...)
	return inst
}
