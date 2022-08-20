package storage

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"

	"github.com/victor_diditskiy/replication_experiment/pkg/dbpool"
	"github.com/victor_diditskiy/replication_experiment/pkg/entity"
)

const (
	insertSQL = "insert into data (name, value) values ('{param1}', {param2})"
	updateSQL = `update data set 
				 	name = $1,
					value = $2,
					updated_at = NOW()
				 where id = $3
	`
	getSql   = "select id, name, value, created_at, updated_at from data where id = $1"
	countSql = "select count(*) from data"
)

type CombinedStorage interface {
	Leader
	Follower
}

type Leader interface {
	Insert(data ...entity.Data) error
	Update(data entity.Data) error
}

type Follower interface {
	Get(id int64) (*entity.Data, error)
	Count() (int64, error)
}

type Storage struct {
	dbPool *dbpool.DBPool
}

func New(dbPool *dbpool.DBPool) *Storage {
	return &Storage{
		dbPool: dbPool,
	}
}

func (s *Storage) Insert(data ...entity.Data) error {
	if len(data) == 0 {
		return errors.New("no data passed")
	}
	db := s.dbPool.GetRandomLeader()

	var sql = ""

	for _, item := range data {
		itemSql := strings.Replace(insertSQL, "{param1}", item.Name, 1)
		itemSql = strings.Replace(itemSql, "{param2}", fmt.Sprintf("%d", item.Value), 1)

		sql += itemSql + "; "
	}
	_, err := db.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to insert data")
	}

	return nil
}

func (s *Storage) Update(data entity.Data) error {
	if data.ID == 0 {
		return errors.New("no id set for updating data")
	}

	db := s.dbPool.GetRandomLeader()

	_, err := db.Exec(updateSQL, data.Name, data.Value, data.ID)
	if err != nil {
		return errors.Wrap(err, "failed to update data")
	}

	return nil
}

func (s *Storage) Get(id int64) (*entity.Data, error) {
	follower := s.dbPool.GetRandomFollower()

	rows, err := follower.Query(getSql, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed execute get query")
	}
	defer rows.Close()

	rows.Next()

	data := &entity.Data{}
	err = rows.Scan(&data.ID, &data.Name, &data.Value, &data.CreatedAt, &data.UpdatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "failed scan get query result")
	}

	return data, nil
}

func (s *Storage) Count() (int64, error) {
	follower := s.dbPool.GetRandomFollower()

	rows, err := follower.Query(countSql)
	if err != nil {
		return 0, errors.Wrap(err, "failed execute count query")
	}
	defer rows.Close()

	rows.Next()

	var count int64
	err = rows.Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "failed scan count query result")
	}

	return count, nil
}
