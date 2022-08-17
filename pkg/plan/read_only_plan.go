package plan

import (
	"context"

	"github.com/pkg/errors"

	"github.com/victor_diditskiy/replication_experiment/pkg/workload"
)

type ReadOnlyPlan struct {
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

	workloadConfig := workload.Config{
		ScaleFactor: config.ReadWorkload.ScaleFactor,
	}
	err := p.manager.workloads.StartWorkload(ctx, workload.ReadWorkloadName, workloadConfig)
	if err != nil {
		return errors.Wrap(err, "read-only plan starting failed")
	}

	return nil
}

func (p *ReadOnlyPlan) Stop() error {
	err := p.manager.workloads.StopWorkload(workload.ReadWorkloadName)
	if err != nil {
		return errors.Wrap(err, "read-only plan stopping failed")
	}

	return nil
}
