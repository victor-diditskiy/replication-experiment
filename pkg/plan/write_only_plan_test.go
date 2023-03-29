package plan

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/victor_diditskiy/replication_experiment/pkg/workload/mock"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/victor_diditskiy/replication_experiment/pkg/workload"
)

func TestWriteOnlyPlan_StartStart(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Run("failed to start plan", func(t *testing.T) {
		t.Run("no config passed", func(t *testing.T) {
			workloads := make(workload.Workloads)
			plan := NewWriteOnlyPlan(workloads)

			config := Config{}
			err := plan.Start(config)
			require.EqualError(t, err, "neither insert or update workload config have been set for write-only plan")
		})

		t.Run("write-workload isn't configured", func(t *testing.T) {
			workloads := make(workload.Workloads)
			plan := NewWriteOnlyPlan(workloads)

			config := Config{InsertWorkload: &ConfigItem{}}
			err := plan.Start(config)

			require.EqualError(t, err, "write-only plan starting failed: failed to find insert workload")
		})

		t.Run("failed to start insert workload", func(t *testing.T) {
			workloads := make(workload.Workloads)
			insertWorkload := mock.NewMockWorkload(ctrl)
			workloads[workload.InsertWorkloadName] = insertWorkload
			plan := NewReadOnlyPlan(workloads)

			expectedConfig := workload.Config{
				ScaleFactor: 10,
				MaxItems:    100_000_000,
			}
			insertWorkload.
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
		writeOnlyWorkload := mock.NewMockWorkload(ctrl)
		workloads[workload.ReadWorkloadName] = writeOnlyWorkload
		plan := NewReadOnlyPlan(workloads)

		expectedConfig := workload.Config{
			ScaleFactor: 10,
			MaxItems:    100_000_000,
		}
		writeOnlyWorkload.
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
