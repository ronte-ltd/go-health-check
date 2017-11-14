package checkers

import (
	"sync"
)

// Map it is thread-safe map implementation
type Map struct {
	mx sync.Mutex
	m  map[string]Checker
}

//NewCheckersMap return instance CheckersMap without initial capacity
func NewCheckersMap() *Map {
	return &Map{
		m: make(map[string]Checker),
	}
}

//NewCheckersMapWithLen return instance CheckersMap with initial capacity
func NewCheckersMapWithLen(length int) *Map {
	return &Map{
		m: make(map[string]Checker, length),
	}
}

//Store add new pair in map
func (c *Map) Store(key string, value Checker) {
	c.mx.Lock()
	c.m[key] = value
	c.mx.Unlock()
}

//Load value by key from map
func (c *Map) Load(key string) (Checker, bool) {
	c.mx.Lock()
	health, ok := c.m[key]
	c.mx.Unlock()
	return health, ok
}

//Range applay input function on each pair in map, if `func` return false, iteration break
func (c *Map) Range(f func(key string, value Checker) bool) {
	c.mx.Lock()
	for k, v := range c.m {
		if !f(k, v) {
			break
		}
	}
	c.mx.Unlock()
}

//Len return count items in map
func (c *Map) Len() int {
	c.mx.Lock()
	length := len(c.m)
	c.mx.Unlock()
	return length
}
