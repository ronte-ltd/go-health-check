package checkers

import (
	"github.com/rainycape/memcache"
	"testing"
)

func IgnoreNewMemcachedChecker(t *testing.T) {
	mc, err := memcache.New("127.0.0.1:11211")
	defer mc.Close()
	if err != nil {
		t.Fatalf("Cannot connect to server: %s", err.Error())
	}
	mcChecker := NewMemcachedChecker("memcached", mc)
	health, err := mcChecker.Check()
	if err != nil {
		t.Fatalf("Unhealthy: %s", err.Error())
	}
	t.Logf("Health: %+v", health)
}
