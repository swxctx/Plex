package plex

import (
	"sync"
)

// connStore
type connStore struct {
	// client data
	data map[string]*storeInfo
	// lock
	rwMu sync.RWMutex
}

// newConnStore
func newConnStore(capacity ...int) *connStore {
	var (
		cap int
	)
	if len(capacity) > 0 {
		cap = capacity[0]
	}
	return &connStore{
		data: make(map[string]*storeInfo, cap),
	}
}

// get
func (cs *connStore) get(key string) (value *storeInfo, exists bool) {
	cs.rwMu.RLock()
	value, exists = cs.data[key]
	cs.rwMu.RUnlock()
	return value, exists
}

// set
func (cs *connStore) set(key string, value *storeInfo) {
	cs.rwMu.Lock()
	cs.data[key] = value
	cs.rwMu.Unlock()
}

// getOrSet if
func (cs *connStore) getOrSet(key string, value *storeInfo) (actual *storeInfo, exists bool) {
	cs.rwMu.Lock()
	actual, exists = cs.data[key]
	if !exists {
		cs.data[key] = value
		actual = value
	}
	cs.rwMu.Unlock()
	return actual, exists
}

// del
func (cs *connStore) del(key string) {
	cs.rwMu.Lock()
	delete(cs.data, key)
	cs.rwMu.Unlock()
}

// clear clears all data
func (cs *connStore) clear() {
	cs.rwMu.Lock()
	for k := range cs.data {
		delete(cs.data, k)
	}
	cs.rwMu.Unlock()
}

// len return data length
func (cs *connStore) len() int {
	cs.rwMu.RLock()
	defer cs.rwMu.RUnlock()
	return len(cs.data)
}
