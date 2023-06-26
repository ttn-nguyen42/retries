package retries

import (
	"context"
	"errors"
	"time"
)

var (
	ErrFinsihed = errors.New("finished all attempts")
)

type BreakableTask func() error

func Do(task BreakableTask, opts ...Option) error {
	c := newDefaultConfig()
	for _, o := range opts {
		o(&c)
	}
	err := c.ctx.Err()
	if err != nil {
		return err
	}
	var trials uint = 0
	timer := time.NewTimer(0)
	// try until succeed
	if c.attempts <= 0 {
		for {
			err := task()
			if err == nil {
				return nil
			}
			if !c.retryIf(err) {
				return err
			}
			trials += 1
			d := c.delayType(trials, err, &c)
			c.onRetry(trials, d, err)
			timer.Reset(d)
			select {
			case <-timer.C:
			case <-c.ctx.Done():
				timer.Stop()
				return c.ctx.Err()
			}
		}
	}
	for {
		if (c.attempts - 1) == trials {
			return ErrFinsihed
		}
		err := task()
		if err == nil {
			return nil
		}
		if !c.retryIf(err) {
			return err
		}
		trials += 1
		d := c.delayType(trials, err, &c)
		c.onRetry(trials, d, err)
		timer.Reset(d)
		select {
		case <-timer.C:
		case <-c.ctx.Done():
			timer.Stop()
		}
	}
}

func newDefaultConfig() Config {
	return Config{
		attempts: 3,
		delay:    time.Millisecond * 100,
		onRetry:  func(trial uint, incomingTimeout time.Duration, err error) {},
		retryIf: func(err error) bool {
			return true
		},
		delayType:  BackoffDelay,
		maxBackoff: 0,
		ctx:        context.Background(),
	}
}
