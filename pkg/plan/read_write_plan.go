package plan

import (
	"context"

	"github.com/pkg/errors"

	"github.com/victor_diditskiy/replication_experiment/pkg/workload"
)

type ReadWritePlan struct {
	manager         *Manager
	activeWorkloads []workload.Name
}

func NewReadWritePlan(manager *Manager) *ReadWritePlan {
	return &ReadWritePlan{
		manager: manager,
	}
}

func (p *ReadWritePlan) Start(config Config) error {
	if config.InsertWorkload == nil &&
		config.UpdateWorkload == nil &&
		config.ReadWorkload == nil {
		return errors.New("neither read, insert or update workload config have been set for read-xwrite plan")
	}

	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, defaultWorkloadTTL)

	if config.InsertWorkload != nil {
		insertConfig := workload.Config{
			ScaleFactor: config.InsertWorkload.ScaleFactor,
			BatchSize:   config.ReadWorkload.BatchSize,
		}
		err := p.manager.workloads.StartWorkload(ctx, workload.InsertWorkloadName, insertConfig)
		if err != nil {
			return errors.Wrap(err, "read-write plan starting failed")
		}

		p.activeWorkloads = append(p.activeWorkloads, workload.InsertWorkloadName)
	}

	if config.UpdateWorkload != nil {
		updateConfig := workload.Config{
			ScaleFactor: config.UpdateWorkload.ScaleFactor,
			MaxItems:    config.ReadWorkload.MaxItems,
		}
		err := p.manager.workloads.StartWorkload(ctx, workload.UpdateWorkloadName, updateConfig)
		if err != nil {
			return errors.Wrap(err, "read-write plan starting failed")
		}

		p.activeWorkloads = append(p.activeWorkloads, workload.UpdateWorkloadName)
	}

	if config.ReadWorkload != nil {
		readConfig := workload.Config{
			ScaleFactor: config.ReadWorkload.ScaleFactor,
			MaxItems:    config.ReadWorkload.MaxItems,
		}
		err := p.manager.workloads.StartWorkload(ctx, workload.ReadWorkloadName, readConfig)
		if err != nil {
			return errors.Wrap(err, "read-write plan starting failed")
		}

		p.activeWorkloads = append(p.activeWorkloads, workload.ReadWorkloadName)
	}

	return nil
}

func (p *ReadWritePlan) Stop() error {
	for _, activeWorkload := range p.activeWorkloads {
		err := p.manager.workloads.StopWorkload(activeWorkload)
		if err != nil {
			return errors.Wrap(err, "read-write plan stopping failed")
		}
	}

	p.activeWorkloads = nil

	return nil
}
