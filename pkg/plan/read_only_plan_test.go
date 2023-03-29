package plan

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/victor_diditskiy/replication_experiment/pkg/workload/mock"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/victor_diditskiy/replication_experiment/pkg/workload"
)

func TestReadOnlyPlan_Start(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Run("failed to start plan", func(t *testing.T) {
		t.Run("no config passed", func(t *testing.T) {
			workloads := make(workload.Workloads)
			plan := NewReadOnlyPlan(workloads)

			config := Config{}
			err := plan.Start(config)
			require.EqualError(t, err, "no read workload config has been set for read only plan")
		})

		t.Run("read-workload isn't configured", func(t *testing.T) {
			workloads := make(workload.Workloads)
			plan := NewReadOnlyPlan(workloads)

			config := Config{ReadWorkload: &ConfigItem{}}
			err := plan.Start(config)

			require.EqualError(t, err, "failed to start read-only plan: failed to find read workload")
		})

		t.Run("failed to start workload", func(t *testing.T) {
			workloads := make(workload.Workloads)
			readOnlyWorkload := mock.NewMockWorkload(ctrl)
			workloads[workload.ReadWorkloadName] = readOnlyWorkload
			plan := NewReadOnlyPlan(workloads)

			expectedConfig := workload.Config{
				ScaleFactor: 10,
				MaxItems:    100_000_000,
			}
			readOnlyWorkload.
				EXPECT().
				Start(gomock.Any(), expectedConfig).
				Return(errors.New("error"))

			config := Config{ReadWorkload: &ConfigItem{
				ScaleFactor: 10,
				MaxItems:    100_000_000,
			}}
			err := plan.Start(config)
			require.EqualError(t, err, "read-only plan starting failed: error")
		})
	})

	t.Run("start plan", func(t *testing.T) {
		workloads := make(workload.Workloads)
		readOnlyWorkload := mock.NewMockWorkload(ctrl)
		workloads[workload.ReadWorkloadName] = readOnlyWorkload
		plan := NewReadOnlyPlan(workloads)

		expectedConfig := workload.Config{
			ScaleFactor: 10,
			MaxItems:    100_000_000,
		}
		readOnlyWorkload.
			EXPECT().
			Start(gomock.Any(), expectedConfig).
			Return(nil)

		config := Config{ReadWorkload: &ConfigItem{
			ScaleFactor: 10,
			MaxItems:    100_000_000,
		}}
		err := plan.Start(config)
		require.NoError(t, err)
	})
}

func TestReadOnlyPlan_Stop(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Run("failed to stop plan", func(t *testing.T) {
		t.Run("no active plan", func(t *testing.T) {
			plan := NewReadOnlyPlan(nil)

			err := plan.Stop()
			require.EqualError(t, err, "no active read-only plan")
		})

		t.Run("failed to stop plan", func(t *testing.T) {
			workloads := make(workload.Workloads)
			readOnlyWorkload := mock.NewMockWorkload(ctrl)
			workloads[workload.ReadWorkloadName] = readOnlyWorkload
			plan := NewReadOnlyPlan(workloads)
			expectedConfig := workload.Config{
				ScaleFactor: 10,
				MaxItems:    100_000_000,
			}
			readOnlyWorkload.
				EXPECT().
				Start(gomock.Any(), expectedConfig).
				Return(nil)

			readOnlyWorkload.
				EXPECT().
				Stop(workload.Config{}).
				Return(errors.New("error"))

			config := Config{ReadWorkload: &ConfigItem{
				ScaleFactor: 10,
				MaxItems:    100_000_000,
			}}
			err := plan.Start(config)
			require.NoError(t, err)

			err = plan.Stop()
			require.EqualError(t, err, "read-only plan stopping failed: error")
		})
	})

	t.Run("start plan", func(t *testing.T) {
		workloads := make(workload.Workloads)
		readOnlyWorkload := mock.NewMockWorkload(ctrl)
		workloads[workload.ReadWorkloadName] = readOnlyWorkload
		plan := NewReadOnlyPlan(workloads)

		expectedConfig := workload.Config{
			ScaleFactor: 10,
			MaxItems:    100_000_000,
		}
		readOnlyWorkload.
			EXPECT().
			Start(gomock.Any(), expectedConfig).
			Return(nil)

		readOnlyWorkload.
			EXPECT().
			Stop(workload.Config{}).
			Return(nil)

		config := Config{ReadWorkload: &ConfigItem{
			ScaleFactor: 10,
			MaxItems:    100_000_000,
		}}
		err := plan.Start(config)
		require.NoError(t, err)

		err = plan.Stop()
		require.NoError(t, err)
	})
}
