package barrier

import (
	"sync"
)

type Barrier struct {
	trip     int
	counter  int
	MuxBarr  *sync.Mutex
	WaitBarr *sync.WaitGroup
}

func NewBarrier(n int) *Barrier {
	MutexBarrier := &sync.Mutex{}
	WaitBarrier := &sync.WaitGroup{}
	WaitBarrier.Add(1)
	return &Barrier{n, 0, MutexBarrier, WaitBarrier}
}

func (b *Barrier) Wait() {
	b.MuxBarr.Lock()
	prBarr := b.WaitBarr //nombre
	b.counter++
	if b.counter < b.trip {
		b.MuxBarr.Unlock()
		prBarr.Wait()
	} else if b.counter == b.trip{
		b.WaitBarr.Done()
		b.WaitBarr = &sync.WaitGroup{}
		b.WaitBarr.Add(1)
		b.counter = 0
		b.MuxBarr.Unlock()
	}
}
