package plan

import (
	"context"

	"github.com/pkg/errors"

	"github.com/victor_diditskiy/replication_experiment/pkg/workload"
)

type WriteOnlyPlan struct {
	manager         *Manager
	activeWorkloads []workload.Name
}

func NewWriteOnlyPlan(manager *Manager) *WriteOnlyPlan {
	return &WriteOnlyPlan{
		manager: manager,
	}
}

func (p *WriteOnlyPlan) Start(config Config) error {
	if config.InsertWorkload == nil && config.UpdateWorkload == nil {
		return errors.New("neither insert or update workload config have been set for write-only plan")
	}

	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, defaultWorkloadTTL)

	if config.InsertWorkload != nil {
		insertConfig := workload.Config{
			ScaleFactor: config.InsertWorkload.ScaleFactor,
			BatchSize:   config.InsertWorkload.BatchSize,
		}
		err := p.manager.workloads.StartWorkload(ctx, workload.InsertWorkloadName, insertConfig)
		if err != nil {
			return errors.Wrap(err, "write-only plan starting failed")
		}

		p.activeWorkloads = append(p.activeWorkloads, workload.InsertWorkloadName)
	}

	if config.UpdateWorkload != nil {
		updateConfig := workload.Config{ScaleFactor: config.UpdateWorkload.ScaleFactor}
		err := p.manager.workloads.StartWorkload(ctx, workload.UpdateWorkloadName, updateConfig)
		if err != nil {
			return errors.Wrap(err, "write-only plan starting failed")
		}

		p.activeWorkloads = append(p.activeWorkloads, workload.UpdateWorkloadName)
	}

	return nil
}

func (p *WriteOnlyPlan) Stop() error {
	for _, activeWorkload := range p.activeWorkloads {
		err := p.manager.workloads.StopWorkload(activeWorkload)
		if err != nil {
			return errors.Wrap(err, "write-only plan stopping failed")
		}
	}

	p.activeWorkloads = nil

	return nil
}
