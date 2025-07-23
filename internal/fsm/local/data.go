package local

import "sync"

type data struct {
	state  string
	values map[string]string
	rm     sync.RWMutex
}

func newData(state string) *data {
	return &data{
		state:  state,
		values: make(map[string]string),
	}
}

func (d *data) State() string {
	d.rm.RLock()
	defer d.rm.RUnlock()
	return d.state
}

func (d *data) setState(state string) {
	d.rm.Lock()
	defer d.rm.Unlock()
	d.state = state
}

func (d *data) GetValue(param string) (string, bool) {
	d.rm.RLock()
	defer d.rm.RUnlock()
	value, ok := d.values[param]
	return value, ok
}

func (d *data) SetValue(param string, value string) bool {
	d.rm.Lock()
	defer d.rm.Unlock()

	d.values[param] = value
	return true
}
