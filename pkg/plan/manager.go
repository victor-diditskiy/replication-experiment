package plan

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"

	"github.com/victor_diditskiy/replication_experiment/pkg/workload"
)

const (
	StartCommand CommandName = "start"
	StopCommand  CommandName = "stop"
)

type CommandName string

type Plans map[Name]Plan

type Manager struct {
	hasActivePlan bool
	m             sync.Mutex

	activePlan Plan
	plans      map[Name]Plan
	workloads  workload.Workloads
}

type Command struct {
	Name     CommandName
	PlanName Name
	Config   Config
}

func NewManager(workloads workload.Workloads) *Manager {
	manager := &Manager{
		m: sync.Mutex{},

		workloads: workloads,
	}
	plans := Plans{
		ReadOnlyPlanName:  NewReadOnlyPlan(manager),
		WriteOnlyPlanName: NewWriteOnlyPlan(manager),
		ReadWritePlanName: NewReadWritePlan(manager),
	}
	manager.plans = plans

	return manager
}

func (m *Manager) Execute(command Command) error {
	m.m.Lock()
	defer m.m.Unlock()

	var err error
	switch command.Name {
	case StartCommand:
		plan, ok := m.plans[command.PlanName]
		if !ok {
			return fmt.Errorf("undefined plan name \"%s\" given", command.PlanName)
		}

		err = m.startPlan(command.Config, plan)
	case StopCommand:
		err = m.stopPlan()
	default:
		return fmt.Errorf("got undefined command \"%s\"", command.Name)
	}

	if err != nil {
		return errors.Wrapf(err, "failed to execute \"%s\" command", command.Name)
	}

	return nil
}

func (m *Manager) startPlan(planConfig Config, plan Plan) error {
	if m.activePlan != nil {
		return errors.New("another plan has already been activated")
	}

	err := plan.Start(planConfig)
	if err != nil {
		return errors.Wrap(err, "failed to start plan")
	}

	m.activePlan = plan

	return nil
}

func (m *Manager) stopPlan() error {
	if m.activePlan == nil {
		return errors.New("no active plan to stop")
	}

	err := m.activePlan.Stop()
	if err != nil {
		return errors.Wrap(err, "failed to stop plan")
	}

	m.activePlan = nil

	return nil
}
