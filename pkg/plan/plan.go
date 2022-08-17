package plan

import (
	"time"
)

const (
	defaultWorkloadTTL = 1 * time.Hour

	ReadOnlyPlanName  Name = "read-only"
	WriteOnlyPlanName Name = "write-only"
	ReadWritePlanName Name = "read-write"
)

type Name string

type ConfigItem struct {
	ScaleFactor int
}
type Config struct {
	ReadWorkload   *ConfigItem
	InsertWorkload *ConfigItem
	UpdateWorkload *ConfigItem
}

type Plan interface {
	Start(config Config) error
	Stop() error
}
