package smap

import (
	"sync"
)

type STMap struct {
	sync.RWMutex
	mp map[int64]struct{}
}

func (m *STMap) Set(k int64) {
	defer m.Unlock()
	m.Lock()
	m.mp[k] = struct{}{}
}

func (m *STMap) Get(k int64) bool {
	defer m.RUnlock()
	m.RLock()
	if _, ok := m.mp[k]; ok {
		return ok
	}
	return false
}

func (m *STMap) GetAllKeys() []int64 {
	defer m.RUnlock()
	m.RLock()
	res := make([]int64, 0, len(m.mp))
	for k, _ := range m.mp {
		res = append(res, k)
	}
	return res
} 

func NewSTMap(cap int) *STMap{
	m := &STMap{
		mp: make(map[int64]struct{}, cap),
	}
	return m
}