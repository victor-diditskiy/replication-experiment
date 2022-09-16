package workload

import (
	"context"
	"fmt"

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
	BatchSize   int
	MaxItems    int
}

func (w *Workloads) StartWorkload(
	ctx context.Context,
	name Name,
	config Config,
) error {
	wl, err := w.findByName(name)
	if err != nil {
		return err
	}

	err = wl.Start(ctx, config)
	if err != nil {
		return errors.Wrapf(err, "failed to start %s workload", name)
	}

	return nil
}

func (w *Workloads) StopWorkload(name Name) error {
	wl, err := w.findByName(name)
	if err != nil {
		return err
	}

	err = wl.Stop(Config{})
	if err != nil {
		return errors.Wrapf(err, "failed to stop %s workload", name)
	}

	return nil
}

func (w *Workloads) findByName(name Name) (Workload, error) {
	wl, ok := (*w)[name]
	if !ok {
		return nil, fmt.Errorf("failed to find %s workload", name)
	}

	return wl, nil
}

func (c *Config) Validate() error {
	if c.ScaleFactor < 1 {
		return errors.New("scale factor can not be less then 1")
	}

	return nil
}
