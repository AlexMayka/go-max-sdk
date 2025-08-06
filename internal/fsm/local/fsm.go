package local

import (
	"github.com/AlexMayka/go-max-sdk/internal/core"
	"sync"
)

type FSM struct {
	fsm map[int64]*data
	rm  sync.RWMutex
}

func NewFSM() core.FSM {
	return &FSM{fsm: make(map[int64]*data), rm: sync.RWMutex{}}
}

func (f *FSM) CreateUser(id int64, state string) bool {
	f.rm.Lock()
	defer f.rm.Unlock()

	if _, ok := f.fsm[id]; ok {
		return false
	}

	f.fsm[id] = newData(state)
	return true
}

func (f *FSM) GetStateUser(id int64) (string, bool) {
	f.rm.RLock()
	defer f.rm.RUnlock()

	if value, ok := f.fsm[id]; ok {
		return value.State(), true
	}

	return "", false
}

func (f *FSM) SetStateUser(id int64, state string) bool {
	f.rm.Lock()
	defer f.rm.Unlock()

	if value, ok := f.fsm[id]; ok {
		value.setState(state)
		return true
	}

	return false
}

func (f *FSM) SetValue(id int64, param string, data string) bool {
	f.rm.RLock()
	value, ok := f.fsm[id]
	f.rm.RUnlock()

	if ok {
		return value.SetValue(param, data)
	}

	return false
}

func (f *FSM) GetValue(id int64, param string) (string, bool) {
	f.rm.RLock()
	value, ok := f.fsm[id]
	f.rm.RUnlock()

	if ok {
		return value.GetValue(param)
	}

	return "", false
}
