package plan

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/victor_diditskiy/replication_experiment/pkg/workload"
)

type ReadOnlyPlan struct {
	workloads      workload.Workloads
	activeWorkload workload.Workload
}

func NewReadOnlyPlan(workLoads workload.Workloads) *ReadOnlyPlan {
	return &ReadOnlyPlan{
		workloads: workLoads,
	}
}

func (p *ReadOnlyPlan) Start(config Config) error {
	if config.ReadWorkload == nil {
		return errors.New("no read workload config has been set for read only plan")
	}

	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, defaultWorkloadTTL)

	wl, err := p.findByName(workload.ReadWorkloadName)
	if err != nil {
		return errors.Wrap(err, "failed to start read-only plan")
	}

	workloadConfig := workload.Config{
		ScaleFactor: config.ReadWorkload.ScaleFactor,
		MaxItems:    config.ReadWorkload.MaxItems,
	}
	err = wl.Start(ctx, workloadConfig)
	if err != nil {
		return errors.Wrap(err, "read-only plan starting failed")
	}

	p.activeWorkload = wl

	return nil
}

func (p *ReadOnlyPlan) Stop() error {
	if p.activeWorkload == nil {
		return errors.New("no active read-only plan")
	}

	err := p.activeWorkload.Stop(workload.Config{})
	if err != nil {
		return errors.Wrap(err, "read-only plan stopping failed")
	}

	return nil
}

func (p *ReadOnlyPlan) findByName(name workload.Name) (workload.Workload, error) {
	wl, ok := p.workloads[name]
	if !ok {
		return nil, fmt.Errorf("failed to find %s workload", name)
	}

	return wl, nil
}
