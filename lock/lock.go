package lock

import "sync"

type lock struct {
	lock sync.RWMutex
}

func (l *lock) Lock() {
	l.lock.Lock()
}

func (l *lock) Unlock() {
	l.lock.Unlock()
}

func (l *lock) RLock() {
	l.lock.RLock()
}

func (l *lock) RUnlock() {
	l.lock.RUnlock()
}
