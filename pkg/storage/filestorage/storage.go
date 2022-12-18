package filestorage

import (
	"context"
	"os"
	"sync"
)

type Storage struct {
	dir string
	wg  *sync.WaitGroup
}

func NewStorage(dirPath string) *Storage {
	_ = os.Mkdir(dirPath, os.ModePerm)
	return &Storage{dir: dirPath, wg: &sync.WaitGroup{}}
}

func (s *Storage) Close(ctx context.Context) chan struct{} {
	done := make(chan struct{})

	doneWaiting := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(doneWaiting)
	}()

	go func() {
		defer close(done)
		select { // either waits for the messages to process or timeout from context
		case <-doneWaiting:
		case <-ctx.Done():
		}
	}()

	return done
}
