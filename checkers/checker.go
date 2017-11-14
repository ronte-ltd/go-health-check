package checkers

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/rainycape/memcache"
	mgo "gopkg.in/mgo.v2"
)

const addedNewChecker = "Added new checker '%s' type of %s checker"

// Checker provide common method health checking
type Checker interface {
	Check() (Health, error)
	Name() string
}

//HealthChecker struct provide Health struct for current checker and arrya all sub-checker
type HealthChecker struct {
	Health
	Checkers *Map `json:"checkers,omitempty"`
	handler  Handler
	logging  chan string
}

//Check iterate through all subchecker, collect all health and return top-level Health or error
func (hc *HealthChecker) Check() (Health, error) {
	if hc.Checkers.Len() == 0 {
		hc.Msg = "len checkers is 0"
		hc.PushMessage(hc.Health.ToString())
		return hc.Health, nil
	}

	if hc.SubHealth == nil {
		hc.SubHealth = NewSubHealthMapWithLen(hc.Checkers.Len())
	}
	hc.Up()

	wg := sync.WaitGroup{}
	wg.Add(hc.Checkers.Len())

	f := func(key string, value Checker) bool {
		h, err := value.Check()
		if err != nil {
			h = HealthError(err)
		}
		if hc.Status == "UP" && h.Status == "DOWN" {
			hc.Down()
		}
		hc.AddSubHealth(value.Name(), h)
		wg.Done()
		return true
	}

	go hc.Checkers.Range(f)

	wg.Wait()
	hc.PushMessage(hc.Health.ToString())
	return hc.Health, nil
}

//Name return name current health check
func (hc *HealthChecker) Name() string {
	return hc.Health.Name
}

//Registry added new checker as child for current health checker
func (hc *HealthChecker) Registry(name string, checker Checker) {
	hc.Checkers.Store(name, checker)
	hc.PushMessage(fmt.Sprintf(addedNewChecker, name, "composite"))
}

//RegistryURL added new HTTP checker as child by name
func (hc *HealthChecker) RegistryURL(name, url string) {
	hc.Checkers.Store(name, NewHTTPChecker(name, url))
	hc.PushMessage(fmt.Sprintf(addedNewChecker, name, "HTTP"))
}

//RegistryFunc added new `func` checker as child by name
func (hc *HealthChecker) RegistryFunc(name string, checkFunc func() Health) {
	hc.Checkers.Store(name, NewFuncChecker(name, checkFunc))
	hc.PushMessage(fmt.Sprintf(addedNewChecker, name, "`func`"))
}

//RegistryDB added new DB checker as child by name
func (hc *HealthChecker) RegistryDB(name string, db *sql.DB) {
	hc.Checkers.Store(name, NewDBChecker(name, db))
	hc.PushMessage(fmt.Sprintf(addedNewChecker, name, "SQL DB"))
}

//RegistryMemcached added new memcached checker as child by name
func (hc *HealthChecker) RegistryMemcached(name string, client *memcache.Client) {
	hc.Checkers.Store(name, NewMemcachedChecker(name, client))
	hc.PushMessage(fmt.Sprintf(addedNewChecker, name, "Memcached"))
}

//RegistryMongo added new Mongodb checker as child by name
func (hc *HealthChecker) RegistryMongo(name string, session *mgo.Session) {
	hc.Checkers.Store(name, NewMongoChecker(name, session))
	hc.PushMessage(fmt.Sprintf(addedNewChecker, name, "Mongodb"))
}

//Handle bind address and added routing to this checker
func (hc *HealthChecker) Handle(addr, route string) error {
	hc.handler = NewHandler(hc, addr, route)
	hc.PushMessage(fmt.Sprintf("Start '%s' on '%s' route '%s'", hc.Name(), addr, route))
	return hc.handler.Server.ListenAndServe()
}

// NewHealthChecker return new instance HealthChecker with default parameters with disable logging channel
func NewHealthChecker(name string) HealthChecker {
	return HealthChecker{
		Health: Health{
			Name:   name,
			Status: DOWN,
		},
		Checkers: NewCheckersMap(),
	}
}

// NewHealthCheckerWithLogger return new instance HealthChecker with default parameters with enable logging channel
func NewHealthCheckerWithLogger(name string) HealthChecker {
	return HealthChecker{
		Health: Health{
			Name:   name,
			Status: DOWN,
		},
		Checkers: NewCheckersMap(),
		logging:  make(chan string),
	}
}

// Up toggle status current checker to `UP`
func (hc *HealthChecker) Up() {
	hc.Status = UP
}

// Down toggle status current checker to `DOWN`
func (hc *HealthChecker) Down() {
	hc.Status = DOWN
}

// DownError toggle status current checker to `DOWN` and set error message to message this checker
func (hc *HealthChecker) DownError(err error) {
	hc.Status = DOWN
	hc.Msg = err.Error()
}

// AddSubHealth added Health of the chekcer as children to current checker via name
func (hc *HealthChecker) AddSubHealth(name string, health Health) {
	hc.SubHealth.Store(name, health)
}

// GetLogger return chan string that you can wrapped in your logger format
func (hc *HealthChecker) GetLogger() <-chan string {
	return hc.logging
}

// PushMessage added new message in the chan logger
func (hc *HealthChecker) PushMessage(msg string) {
	if hc.logging != nil {
		hc.logging <- msg
	}
}

// PushHealth send current health to logging channel
func (hc *HealthChecker) PushHealth() {
	if hc.logging != nil {
		hc.logging <- hc.Health.ToString()
	}
}

// CheckError send error to loggin channel
func (hc *HealthChecker) CheckError(err error) {
	if err != nil {
		msg := fmt.Sprintf("Was error: %s", err.Error())
		hc.PushMessage(msg)
	}
}
