package local

import (
	"sync"

	"github.com/nimbolus/terraform-backend/pkg/terraform"
)

type Lock struct {
	mutex sync.Mutex
	db    map[string][]byte
}

func NewLock() *Lock {
	return &Lock{
		db: make(map[string][]byte),
	}
}

func (l *Lock) GetName() string {
	return "local"
}

func (l *Lock) Lock(s *terraform.State) (bool, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	lock, ok := l.db[s.ID]
	if ok {
		if string(lock) == string(s.Lock) {
			// you already have the lock
			return true, nil
		}

		s.Lock = lock

		return false, nil
	}

	l.db[s.ID] = s.Lock

	return true, nil
}

func (l *Lock) Unlock(s *terraform.State) (bool, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	lock, ok := l.db[s.ID]
	if !ok {
		return false, nil
	}

	if string(lock) != string(s.Lock) {
		s.Lock = lock

		return false, nil
	}

	delete(l.db, s.ID)

	return true, nil
}
