package storage

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/victor_diditskiy/replication_experiment/pkg/entity"
	"github.com/victor_diditskiy/replication_experiment/pkg/metric"
)

type MetricStorage struct {
	storage CombinedStorage
}

func NewMetricStorage(storage CombinedStorage) *MetricStorage {
	return &MetricStorage{
		storage: storage,
	}
}

func (s *MetricStorage) Insert(data entity.Data) error {
	t := time.Now()
	defer metric.RequestCounter.With(prometheus.Labels{metric.TypeLabel: metric.WriteType, metric.OperationLabel: metric.InsertOperation}).Inc()
	defer func() {
		metric.RequestTimingCounter.
			With(prometheus.Labels{metric.TypeLabel: metric.WriteType, metric.OperationLabel: metric.InsertOperation}).
			Observe(float64(time.Since(t).Milliseconds()))
	}()

	return s.storage.Insert(data)
}

func (s *MetricStorage) Update(data entity.Data) error {
	t := time.Now()
	defer metric.RequestCounter.With(prometheus.Labels{metric.TypeLabel: metric.WriteType, metric.OperationLabel: metric.UpdateOperation}).Inc()
	defer func() {
		metric.RequestTimingCounter.
			With(prometheus.Labels{metric.TypeLabel: metric.WriteType, metric.OperationLabel: metric.UpdateOperation}).
			Observe(float64(time.Since(t).Milliseconds()))
	}()

	return s.storage.Update(data)
}

func (s *MetricStorage) Get(id int64) (*entity.Data, error) {
	t := time.Now()
	defer metric.RequestCounter.With(prometheus.Labels{metric.TypeLabel: metric.ReadType, metric.OperationLabel: metric.GetLOperation}).Inc()
	defer func() {
		metric.RequestTimingCounter.
			With(prometheus.Labels{metric.TypeLabel: metric.ReadType, metric.OperationLabel: metric.GetLOperation}).
			Observe(float64(time.Since(t).Milliseconds()))
	}()

	return s.storage.Get(id)

}

func (s *MetricStorage) Count() (int64, error) {
	t := time.Now()
	defer metric.RequestCounter.With(prometheus.Labels{metric.TypeLabel: metric.ReadType, metric.OperationLabel: metric.CountOperation}).Inc()
	defer func() {
		metric.RequestTimingCounter.
			With(prometheus.Labels{metric.TypeLabel: metric.ReadType, metric.OperationLabel: metric.CountOperation}).
			Observe(float64(time.Since(t).Milliseconds()))
	}()

	return s.storage.Count()

}
