package counter

import (
	"context"
	"sync"
)

type Counter struct {
	value     int
	limit     int
	increment chan int
}

func NewCounter(limit int) *Counter {
	c := Counter{
		limit:     limit,
		value:     0,
		increment: make(chan int, limit),
	}
	return &c
}

func (c *Counter) Add(amount int, ctx context.Context, cancel context.CancelFunc) {
	if c.value >= c.limit {
		ctx.Done()
		cancel()
	}
	c.increment <- amount
}

func (c *Counter) Increment(wg *sync.WaitGroup, cancel context.CancelFunc) {
	defer wg.Done()
	for step := range c.increment {
		if c.value < c.limit {
			c.value += step
		} else {
			cancel()
		}
	}
}

func (c *Counter) Value() int {
	return c.value
}

func (c *Counter) CloseChanel() {
	close(c.increment)
}
