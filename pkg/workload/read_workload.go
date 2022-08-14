package workload

import (
	"context"
	"fmt"
	"math/rand"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/victor_diditskiy/replication_experiment/pkg/storage"
)

type ReadWorkload struct {
	internalCtx   context.Context
	contextCancel context.CancelFunc
	m             sync.Mutex

	started bool

	log    logrus.FieldLogger
	reader storage.Follower
}

func NewReadWorkload(log logrus.FieldLogger, reader storage.Follower) *ReadWorkload {
	return &ReadWorkload{
		m:       sync.Mutex{},
		started: false,
		log:     log,
		reader:  reader,
	}
}

func (rw *ReadWorkload) Start(ctx context.Context, conf Config) error {
	rw.m.Lock()
	defer rw.m.Unlock()

	if rw.started {
		return errors.New("read workload has already started")
	}

	if err := conf.Validate(); err != nil {
		return errors.Wrap(err, "invalid workload config")
	}

	ctx, cancel := context.WithCancel(ctx)
	rw.internalCtx = ctx
	rw.contextCancel = cancel

	for i := 0; i < conf.ScaleFactor; i++ {
		go func() {
			for {
				select {
				case <-rw.internalCtx.Done():
					return
				default:
				}

				count, err := rw.reader.Count()
				if err != nil {
					rw.log.
						Error(fmt.Sprintf("failed to count data at storage: %s", err))
				}

				id := rand.Int63n(count) + 1
				_, err = rw.reader.Get(conf.MaxEntries)
				if err != nil {
					rw.log.
						WithField("id", id).
						Error(fmt.Sprintf("failed to get data from storage: %s", err))
				}
			}
		}()
	}

	rw.started = true

	return nil
}

func (rw *ReadWorkload) Stop(_ Config) error {
	rw.m.Lock()
	defer rw.m.Unlock()

	if !rw.started {
		return errors.New("impossible to stop read workload because it has not been started")
	}

	rw.contextCancel()

	return nil
}
