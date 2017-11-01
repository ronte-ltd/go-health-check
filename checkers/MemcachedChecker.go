package checkers

import (
	"errors"
	"github.com/rainycape/memcache"
)

const KEY_STATUS = "SelfHealthCheckStatus"

type MemcachedChecker struct {
	HealthChecker HealthChecker
	McClient      *memcache.Client
}

// Create new MemcachedChecker by client of Memcached
func NewMemcachedChecker(name string, client *memcache.Client) MemcachedChecker {
	return MemcachedChecker{
		HealthChecker: NewHealthChecker(name),
		McClient:      client,
	}
}

func (mc *MemcachedChecker) Check() (Health, error) {
	item := &memcache.Item{Key: KEY_STATUS, Value: []byte("OK")}
	err := mc.McClient.Set(item)
	if err != nil {
		mc.HealthChecker.DownError(err)
		return mc.HealthChecker.Health, err
	}
	res, err := mc.McClient.Get(KEY_STATUS)
	if err != nil {
		mc.HealthChecker.DownError(err)
		return mc.HealthChecker.Health, err
	}
	if string(res.Value) != "OK" {
		err := errors.New("value in Cache will changed")
		mc.HealthChecker.DownError(err)
		return mc.HealthChecker.Health, err
	}
	mc.HealthChecker.Up()
	return mc.HealthChecker.Health, nil
}

func (mc *MemcachedChecker) Name() string {
	return mc.HealthChecker.Name
}
