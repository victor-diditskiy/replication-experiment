package workload

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/victor_diditskiy/replication_experiment/pkg/entity"
	"github.com/victor_diditskiy/replication_experiment/pkg/storage"
)

type InsertWorkload struct {
	internalCtx   context.Context
	contextCancel context.CancelFunc
	m             sync.Mutex

	started bool

	log    logrus.FieldLogger
	writer storage.Leader
}

func NewInsertWorkload(log logrus.FieldLogger, writer storage.Leader) *InsertWorkload {
	return &InsertWorkload{
		m:       sync.Mutex{},
		started: false,
		log:     log,
		writer:  writer,
	}
}

func (iw *InsertWorkload) Start(ctx context.Context, conf Config) error {
	iw.m.Lock()
	defer iw.m.Unlock()

	if iw.started {
		return errors.New("insert workload has already started")
	}

	if err := conf.Validate(); err != nil {
		return errors.Wrap(err, "invalid workload config")
	}

	ctx, cancel := context.WithCancel(ctx)
	iw.internalCtx = ctx
	iw.contextCancel = cancel

	for i := 0; i < conf.ScaleFactor; i++ {
		go func() {
			for {
				select {
				case <-iw.internalCtx.Done():
					iw.log.Info("insert workload finished")
					return
				default:
				}

				err := iw.writer.Insert(entity.RandomData())
				if err != nil {
					iw.log.
						Error(fmt.Sprintf("failed to insert data to storage: %s", err))
				}
			}
		}()
	}

	iw.started = true

	return nil
}

func (iw *InsertWorkload) Stop(_ Config) error {
	iw.m.Lock()
	defer iw.m.Unlock()

	if !iw.started {
		return errors.New("impossible to stop insert workload because it has not been started")
	}

	iw.contextCancel()
	iw.started = false

	return nil
}
