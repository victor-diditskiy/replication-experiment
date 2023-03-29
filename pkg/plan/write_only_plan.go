package plan

import (
	"context"
	"fmt"
	"github.com/pkg/errors"

	"github.com/victor_diditskiy/replication_experiment/pkg/workload"
)

type WriteOnlyPlan struct {
	workloads       workload.Workloads
	activeWorkloads []workload.Name
}

func NewWriteOnlyPlan(workloads workload.Workloads) *WriteOnlyPlan {
	return &WriteOnlyPlan{
		workloads: workloads,
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

		wl, err := p.findByName(workload.InsertWorkloadName)
		if err != nil {
			return errors.Wrap(err, "failed to start write-only plan")
		}

		err = wl.Start(ctx, insertConfig)
		if err != nil {
			return errors.Wrap(err, "write-only plan starting failed")
		}

		p.activeWorkloads = append(p.activeWorkloads, workload.InsertWorkloadName)
	}

	if config.UpdateWorkload != nil {
		updateConfig := workload.Config{
			ScaleFactor: config.UpdateWorkload.ScaleFactor,
			MaxItems:    config.UpdateWorkload.MaxItems,
		}

		wl, err := p.findByName(workload.InsertWorkloadName)
		if err != nil {
			return errors.Wrap(err, "failed to start write-only plan")
		}

		err = wl.Start(ctx, updateConfig)
		if err != nil {
			return errors.Wrap(err, "write-only plan starting failed")
		}

		p.activeWorkloads = append(p.activeWorkloads, workload.UpdateWorkloadName)
	}

	return nil
}

func (p *WriteOnlyPlan) Stop() error {
	for _, activeWorkload := range p.activeWorkloads {
		err := p.workloads.StopWorkload(activeWorkload)
		if err != nil {
			return errors.Wrap(err, "write-only plan stopping failed")
		}
	}

	p.activeWorkloads = nil

	return nil
}

func (p *WriteOnlyPlan) findByName(name workload.Name) (workload.Workload, error) {
	wl, ok := p.workloads[name]
	if !ok {
		return nil, fmt.Errorf("failed to find %s workload", name)
	}

	return wl, nil
}
