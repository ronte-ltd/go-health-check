package checkers

import (
	"fmt"
	"sync"
)

// Down - Stuats Service is unhealthy
// Up - Status Service is healthy
const (
	DOWN = "DOWN"
	UP   = "UP"
)

//Health struct provide Metainfo about current healthy and his sub-chekcer
type Health struct {
	Name      string        `json:"name"`
	Status    string        `json:"status"`
	Msg       string        `json:"msg,omitempty"`
	SubHealth *SubHealthMap `json:"subHealth,omitempty"`
}

// HealthError create new instance `Health` with `DOWN` status by error
func HealthError(err error) Health {
	return Health{
		Status: DOWN,
		Msg:    err.Error(),
	}
}

// ToString convert Health to string
func (h *Health) ToString() string {
	return fmt.Sprintf("%+v", h)
}

// SubHealthMap it is thread-safe map implementation
type SubHealthMap struct {
	mx sync.Mutex
	m  map[string]Health
}

//NewSubHealthMap return instance SubHealthMap without initial capacity
func NewSubHealthMap() *SubHealthMap {
	return &SubHealthMap{
		m: make(map[string]Health),
	}
}

//NewSubHealthMapWithLen return instance SubHealthMap with initial capacity
func NewSubHealthMapWithLen(length int) *SubHealthMap {
	return &SubHealthMap{
		m: make(map[string]Health, length),
	}
}

//Store add new pair in map
func (shm *SubHealthMap) Store(key string, value Health) {
	shm.mx.Lock()
	shm.m[key] = value
	shm.mx.Unlock()
}

//Load value by key from map
func (shm *SubHealthMap) Load(key string) (Health, bool) {
	shm.mx.Lock()
	health, ok := shm.m[key]
	shm.mx.Unlock()
	return health, ok
}

//Range applay input function on each pair in map, if `func` return false, iteration break
func (shm *SubHealthMap) Range(f func(key string, value Health) bool) {
	shm.mx.Lock()
	for k, v := range shm.m {
		if !f(k, v) {
			break
		}
	}
	shm.mx.Unlock()
}
