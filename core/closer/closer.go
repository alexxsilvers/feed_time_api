package closer

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"github.com/pkg/errors"

	"github.com/alexxsilvers/feed_time_api/core/logger"
)

type Closer struct {
	sync.Mutex
	once  sync.Once
	done  chan struct{}
	funcs []func() error
}

func New(sig ...os.Signal) *Closer {
	c := &Closer{
		done: make(chan struct{}),
	}

	if len(sig) > 0 {
		go func() {
			waitCH := make(chan os.Signal, 1)
			signal.Notify(waitCH, sig...)
			<-waitCH
			signal.Stop(waitCH)
			c.CloseAll()
		}()
	}

	return c
}

func (c *Closer) Add(f func() error) {
	c.Lock()
	c.funcs = append(c.funcs, f)
	c.Unlock()
}

func (c *Closer) Wait() {
	select {
	case <-c.done:
	}
}

func (c *Closer) CloseAll() {
	c.once.Do(func() {
		defer close(c.done)

		c.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.Unlock()

		errs := make(chan error, len(funcs))
		for _, f := range funcs {
			go func(func() error) {
				errs <- f()
			}(f)
		}

		for i := 0; i < cap(errs); i++ {
			err := <-errs
			if err != nil {
				logger.Error(context.Background(), errors.Wrap(err, "closer error"))
			}
		}
	})
}
