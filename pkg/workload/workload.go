package workload

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/victor_diditskiy/replication_experiment/pkg/storage"
)

const (
	ReadWorkloadName   Name = "read"
	InsertWorkloadName Name = "insert"
	UpdateWorkloadName Name = "update"
)

type Name string

type Workload interface {
	Start(ctx context.Context, conf Config) error
	Stop(conf Config) error
}

type Workloads map[Name]Workload

func NewWorkloads(log logrus.FieldLogger, storage storage.CombinedStorage) Workloads {
	return Workloads{
		ReadWorkloadName:   NewReadWorkload(log, storage),
		InsertWorkloadName: NewInsertWorkload(log, storage),
		UpdateWorkloadName: NewUpdateWorkload(log, storage),
	}
}

type Config struct {
	ScaleFactor int
}

func (c *Config) Validate() error {
	if c.ScaleFactor < 1 {
		return errors.New("scale factor can not be less then 1")
	}

	return nil
}
