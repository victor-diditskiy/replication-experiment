package dbpool

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
)

const (
	configDB = "CONFIG_DB"
)

type DB struct {
	DSN             string
	PoolConnections int
}
type Config struct {
	Leaders   []DB
	Followers []DB
}

// DBPool is container for db instances. Is has two group if db instances: leaders and follower.
// Leaders are used to perform write only operations. Followers are used to read only operations.
// It's possible to configure same db as leader and follower to perform both operation types.
type DBPool struct {
	leaders   []*sql.DB
	followers []*sql.DB
}

func NewPool(dbConfigPath string) (*DBPool, error) {
	configDBStr := os.Getenv(configDB)
	if configDBStr == "" {
		return nil, fmt.Errorf("%s has not been set", configDB)
	}

	conf := &Config{}
	err := json.Unmarshal([]byte(configDBStr), conf)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse db config")
	}

	pool := &DBPool{}
	for i, leader := range conf.Leaders {
		if leader.DSN == "" {
			return nil, fmt.Errorf("no DSN set for leader #%d", i)
		}

		db, err := createDB(leader)

		if err != nil {
			return nil, errors.Wrap(err, "failed to init DB connection")
		}

		pool.leaders = append(pool.leaders, db)
	}

	for i, follower := range conf.Followers {
		if follower.DSN == "" {
			return nil, fmt.Errorf("no DSN set for follower #%d", i)
		}

		db, err := createDB(follower)

		if err != nil {
			return nil, errors.Wrap(err, "failed to init DB connection")
		}

		pool.followers = append(pool.followers, db)
	}

	return pool, nil
}

func createDB(dbConfig DB) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbConfig.DSN)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "failed to ping db")
	}

	db.SetMaxOpenConns(dbConfig.PoolConnections)
	return db, nil
}

func (p *DBPool) GetLeaders() []*sql.DB {
	return p.leaders
}

func (p *DBPool) GetRandomLeader() *sql.DB {
	// TODO: configure returning random leader
	return p.leaders[0]
}

func (p *DBPool) GetRandomFollower() *sql.DB {
	// TODO: configure returning random follower
	return p.followers[0]
}
