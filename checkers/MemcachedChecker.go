package checkers

import (
	"errors"

	"github.com/rainycape/memcache"
)

const keyStatus = "SelfHealthCheckStatus"

//MemcachedChecker is ctruct provide Helath checker and client for Memcached
type MemcachedChecker struct {
	HealthChecker HealthChecker
	McClient      *memcache.Client
}

// NewMemcachedChecker return by client of Memcached
func NewMemcachedChecker(name string, client *memcache.Client) *MemcachedChecker {
	return &MemcachedChecker{
		HealthChecker: NewHealthChecker(name),
		McClient:      client,
	}
}

//Check connect to memcached and put `keyStatus` to storage, before get `keyStatus` and compare their
func (mc *MemcachedChecker) Check() (Health, error) {
	item := &memcache.Item{Key: keyStatus, Value: []byte("OK")}
	err := mc.McClient.Set(item)
	checkError(err)
	if err != nil {
		mc.HealthChecker.DownError(err)
		mc.HealthChecker.PushHealth()
		return mc.HealthChecker.Health, err
	}
	res, err := mc.McClient.Get(keyStatus)
	if err != nil {
		mc.HealthChecker.DownError(err)
		mc.HealthChecker.PushHealth()
		return mc.HealthChecker.Health, err
	}
	if string(res.Value) != "OK" {
		err := errors.New("value in Cache will changed")
		mc.HealthChecker.DownError(err)
		mc.HealthChecker.PushHealth()
		return mc.HealthChecker.Health, err
	}
	mc.HealthChecker.Up()
	mc.HealthChecker.PushHealth()
	return mc.HealthChecker.Health, nil
}

//Name return name current health check
func (mc *MemcachedChecker) Name() string {
	return mc.HealthChecker.Health.Name
}
