package metric

import (
	"fmt"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace = "replication_experiment"
	subsystem = "application"
	baseName  = "requests"

	TypeLabel = "type"
	ReadType  = "read"
	WriteType = "write"

	OperationLabel  = "operation"
	InsertOperation = "insert"
	UpdateOperation = "update"
	GetLOperation   = "get"
	CountOperation  = "count"
)

var (
	RequestCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      createMetricName("total"),
		Help:      "Count read and write requests",
	}, []string{TypeLabel, OperationLabel})

	RequestTimingCounter = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      createMetricName("timing"),
		Help:      "Count read and write requests",
		Buckets:   []float64{0.1, 1, 5, 10, 25, 50, 100, 200, 500, 1000, 2000, 5000, 10000},
	}, []string{TypeLabel, OperationLabel})
)

func createMetricName(parts ...string) string {
	return fmt.Sprintf("%s_%s", baseName, strings.Join(parts, "_"))
}
