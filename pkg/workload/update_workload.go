package workload

import (
	"context"
	"fmt"
	"math/rand"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/victor_diditskiy/replication_experiment/pkg/entity"
	"github.com/victor_diditskiy/replication_experiment/pkg/storage"
)

type UpdateWorkload struct {
	internalCtx   context.Context
	contextCancel context.CancelFunc
	m             sync.Mutex

	started bool

	log     logrus.FieldLogger
	storage storage.CombinedStorage
}

func NewUpdateWorkload(log logrus.FieldLogger, storage storage.CombinedStorage) *UpdateWorkload {
	return &UpdateWorkload{
		m:       sync.Mutex{},
		started: false,
		log:     log,
		storage: storage,
	}
}

func (uw *UpdateWorkload) Start(ctx context.Context, conf Config) error {
	uw.m.Lock()
	defer uw.m.Unlock()

	if uw.started {
		return errors.New("update workload has already started")
	}

	if err := conf.Validate(); err != nil {
		return errors.Wrap(err, "invalid workload config")
	}

	ctx, cancel := context.WithCancel(ctx)
	uw.internalCtx = ctx
	uw.contextCancel = cancel

	var entriesCount int64
	countChan := make(chan struct{})
	go func() {
		cnt, err := uw.storage.Count()
		if err != nil {
			uw.log.
				Error(fmt.Sprintf("failed to count data at storage: %s", err))
		}
		entriesCount = cnt
		countChan <- struct{}{}

		for {
			select {
			case <-uw.internalCtx.Done():
				return
			default:
			}

			cnt, err := uw.storage.Count()
			if err != nil {
				uw.log.
					Error(fmt.Sprintf("failed to count data at storage: %s", err))
			}

			entriesCount = cnt
		}
	}()

	<-countChan

	for i := 0; i < conf.ScaleFactor; i++ {
		go func() {
			for {
				select {
				case <-uw.internalCtx.Done():
					return
				default:
				}

				id := rand.Int63n(entriesCount)
				data := entity.RandomData()
				data.ID = id
				err := uw.storage.Update(data)
				if err != nil {
					uw.log.
						WithField("id", id).
						Error(fmt.Sprintf("failed to update data at storage: %s", err))
				}
			}
		}()
	}

	uw.started = true

	return nil
}

func (uw *UpdateWorkload) Stop(_ Config) error {
	uw.m.Lock()
	defer uw.m.Unlock()

	if !uw.started {
		return errors.New("impossible to stop update workload because it has not been started")
	}

	uw.contextCancel()

	return nil
}
