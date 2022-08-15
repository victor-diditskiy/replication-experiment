package plan

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/victor_diditskiy/replication_experiment/pkg/workload"
)

const (
	defaultWorkloadTTL = 1 * time.Hour

	ReadOnlyPlanName  Name = "read_only"
	WriteOnlyPlanName Name = "write_only"
	CombinedPlanName  Name = "combined"
)

type Name string

type ConfigItem struct {
	ScaleFactor int
}
type Config struct {
	ReadWorkload   *ConfigItem
	InsertWorkload *ConfigItem
	UpdateWorkload *ConfigItem
}

type Plan interface {
	Start(config Config) error
	Stop() error
}

type ReadOnlyPlan struct {
	ctx context.Context

	manager *Manager
}

func NewReadOnlyPlan(manager *Manager) *ReadOnlyPlan {
	return &ReadOnlyPlan{
		manager: manager,
	}
}

func (p *ReadOnlyPlan) Start(config Config) error {
	if config.ReadWorkload == nil {
		return errors.New("no read workload config has been set for read only plan")
	}

	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, defaultWorkloadTTL)

	readWorkload, ok := p.manager.workloads[workload.ReadWorkloadName]
	if !ok {
		return errors.New("failed to find read workload")
	}

	workloadConfig := workload.Config{
		ScaleFactor: config.ReadWorkload.ScaleFactor,
	}
	err := readWorkload.Start(ctx, workloadConfig)
	if err != nil {
		return errors.Wrap(err, "failed to start read workload")
	}

	return nil
}
func (p *ReadOnlyPlan) Stop() error {
	readWorkload, ok := p.manager.workloads[workload.ReadWorkloadName]
	if !ok {
		return errors.New("failed to find read workload")
	}

	err := readWorkload.Stop(workload.Config{})
	if err != nil {
		return errors.Wrap(err, "failed to stop read workload")
	}

	return nil
}
