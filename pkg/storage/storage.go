package storage

import (
	"github.com/pkg/errors"

	"github.com/victor_diditskiy/replication_experiment/pkg/dbpool"
	"github.com/victor_diditskiy/replication_experiment/pkg/entity"
)

const (
	insertSQL = "insert into data (name, value) values ($1, $2)"
)

type Leader interface {
	Save(data entity.Data) error
}

type Follower interface {
	Get(id int64) *entity.Data
}

type Storage struct {
	dbPool *dbpool.DBPool
}

func New(dbPool *dbpool.DBPool) *Storage {
	return &Storage{
		dbPool: dbPool,
	}
}

func (s *Storage) Save(data entity.Data) error {
	db := s.dbPool.GetLeader()

	_, err := db.Exec(insertSQL, data.Name, data.Value)
	if err != nil {
		return errors.Wrap(err, "failed to insert data")
	}

	return nil
}
func (s *Storage) Get(id int64) *entity.Data {
	// TODO: implement
	return nil
}
