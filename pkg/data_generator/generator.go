package data_generator

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/victor_diditskiy/replication_experiment/pkg/entity"
	"github.com/victor_diditskiy/replication_experiment/pkg/storage"
)

const (
	defaultBatchSize = 10
	workerLimit      = 50
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

func (g *Generator) Generate(ctx context.Context, batchSize, limit, scaleFactor int64) error {
	ch := make(chan struct{})
	defer close(ch)

	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(ctx)

	for i := int64(0); i < scaleFactor; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case <-ch:
					items := make([]entity.Data, 0, batchSize)
					for i := int64(0); i < batchSize; i++ {
						items = append(items, entity.RandomData())
					}

					err := g.storage.Insert(items...)
					if err != nil {
						g.log.Error(err)
					}
				}
			}
		}()
	}

	for i := int64(0); i < limit/batchSize; i++ {
		ch <- struct{}{}
	}

	cancel()

	wg.Wait()

	return nil
}
