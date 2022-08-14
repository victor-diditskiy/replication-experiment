package workload

import (
	"context"

	"github.com/pkg/errors"
)

type Workload interface {
	Start(ctx context.Context, conf Config) error
	Stop(conf Config) error
}

type Config struct {
	ScaleFactor int
	MaxEntries  int64
}

func (c *Config) Validate() error {
	if c.ScaleFactor < 1 {
		return errors.New("scale factor can not be less then 1")
	}

	if c.MaxEntries < 1 {
		return errors.New("max entries can not be less then 1")
	}

	return nil
}
