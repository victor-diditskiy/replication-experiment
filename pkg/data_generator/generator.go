package data_generator

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/victor_diditskiy/replication_experiment/pkg/entity"
	"github.com/victor_diditskiy/replication_experiment/pkg/storage"
)

const (
	workerLimit = 20
)

type Generator struct {
	log     logrus.FieldLogger
	storage storage.Leader
}

func New(log logrus.FieldLogger, storage storage.Leader) *Generator {
	return &Generator{
		log:     log,
		storage: storage,
	}
}

func (g *Generator) Generate(ctx context.Context, limit int64) error {
	ch := make(chan struct{})
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(ctx)

	for i := 0; i < workerLimit; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case <-ch:
					err := g.storage.Save(entity.RandomData())
					if err != nil {
						g.log.Error(err)
					}
				}
			}
		}()
	}

	// TODO: add concurrent saving data
	for i := int64(0); i < limit; i++ {
		ch <- struct{}{}
	}

	cancel()

	wg.Wait()

	return nil
}
