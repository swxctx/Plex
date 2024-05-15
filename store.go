package plex

import (
	"sync"
)

// connStore
type connStore struct {
	data map[string]*storeInfo
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

// Get
func (cs *connStore) Get(key string) (value *storeInfo, exists bool) {
	cs.rwMu.RLock()
	value, exists = cs.data[key]
	cs.rwMu.RUnlock()
	return value, exists
}

// Set
func (cs *connStore) Set(key string, value *storeInfo) {
	cs.rwMu.Lock()
	cs.data[key] = value
	cs.rwMu.Unlock()
}

// GetOrSet if
func (cs *connStore) GetOrSet(key string, value *storeInfo) (actual *storeInfo, exists bool) {
	cs.rwMu.Lock()
	actual, exists = cs.data[key]
	if !exists {
		cs.data[key] = value
		actual = value
	}
	cs.rwMu.Unlock()
	return actual, exists
}

// Del
func (cs *connStore) Del(key string) {
	cs.rwMu.Lock()
	delete(cs.data, key)
	cs.rwMu.Unlock()
}

// Clear clears all data
func (cs *connStore) Clear() {
	cs.rwMu.Lock()
	for k := range cs.data {
		delete(cs.data, k)
	}
	cs.rwMu.Unlock()
}

// Len return data length
func (cs *connStore) Len() int {
	cs.rwMu.RLock()
	defer cs.rwMu.RUnlock()
	return len(cs.data)
}
