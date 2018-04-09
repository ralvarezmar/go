package sems

import (
	"container/list"
	"sync"
)

type Sem struct {
	token   int
	MuxSem  *sync.Mutex
	Waiting *list.List
}

func NewSem(n int) *Sem {
	MutexSem := &sync.Mutex{}
	l := list.New()
	return &Sem{n, MutexSem, l}
}

func (s *Sem) Down() {
	s.MuxSem.Lock()
	if s.token < 0 {
		panic("Semáforo negativo")
	}else	if s.token == 0 {
		waitlocal := &sync.WaitGroup{}
		waitlocal.Add(1)
		s.Waiting.PushFront(waitlocal)
		s.MuxSem.Unlock()
		waitlocal.Wait()
	}else {
		s.token--
		s.MuxSem.Unlock()
	}
}

func (s *Sem) Up() {
	s.MuxSem.Lock()
	if s.Waiting.Len() > 0 {
		waitlocal := s.Waiting.Back()
		s.Waiting.Remove(waitlocal)
		var wg *sync.WaitGroup
		wg = waitlocal.Value.(*sync.WaitGroup)
		wg.Done()
	}else {
		s.token++
	}
	s.MuxSem.Unlock()
}
