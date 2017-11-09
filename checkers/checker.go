package checkers

import (
	"database/sql"

	"github.com/rainycape/memcache"
	mgo "gopkg.in/mgo.v2"
)

// Checker provide common method health checking
type Checker interface {
	Check() (Health, error)
	Name() string
}

//HealthChecker struct provide Health struct for current checker and arrya all sub-checker
type HealthChecker struct {
	Health
	Checkers map[string]Checker `json:"checkers,omitempty"`
	handler  Handler
}

//Health struct provide Metainfo about current healthy and his sub-chekcer
type Health struct {
	Name      string            `json:"name"`
	Status    string            `json:"status"`
	Msg       string            `json:"msg,omitempty"`
	SubHealth map[string]Health `json:"subHealth,omitempty"`
}

// Down - Stuats Service is unhealthy
// Up - Status Service is healthy
const (
	DOWN = "DOWN"
	UP   = "UP"
)

//Check iterate through all subchecker, collect all health and return top-level Health or error
func (hc *HealthChecker) Check() (Health, error) {
	if len(hc.Checkers) == 0 {
		hc.Msg = "len checkers is 0"
		return hc.Health, nil
	}

	if hc.SubHealth == nil {
		hc.SubHealth = make(map[string]Health, len(hc.Checkers))
	}
	hc.Up()

	for _, c := range hc.Checkers {
		var h, err = c.Check()
		if err != nil {
			h = HealthError(err)
		}
		if hc.Status == "UP" && h.Status == "DOWN" {
			hc.Down()
		}
		hc.AddSubHealth(c.Name(), h)
	}

	return hc.Health, nil
}

//Name return name current health check
func (hc *HealthChecker) Name() string {
	return hc.Health.Name
}

//Registry added new checker as child for current health checker
func (hc *HealthChecker) Registry(name string, chekcer Checker) {
	hc.Checkers[name] = chekcer
}

//RegistryURL added new HTTP checker as child by name
func (hc *HealthChecker) RegistryURL(name, url string) {
	hc.Checkers[name] = NewHTTPChecker(name, url)
}

//RegistryFunc added new `func` checker as child by name
func (hc *HealthChecker) RegistryFunc(name string, checkFunc func() Health) {
	hc.Checkers[name] = NewFuncChecker(name, checkFunc)
}

//RegistryDB added new DB checker as child by name
func (hc *HealthChecker) RegistryDB(name string, db *sql.DB) {
	hc.Checkers[name] = NewDBChecker(name, db)
}

//RegistryMemcached added new memcached checker as child by name
func (hc *HealthChecker) RegistryMemcached(name string, client *memcache.Client) {
	hc.Checkers[name] = NewMemcachedChecker(name, client)
}

//RegistryMongo added new Mongodb checker as child by name
func (hc *HealthChecker) RegistryMongo(name string, session *mgo.Session) {
	hc.Checkers[name] = NewMongoChecker(name, session)
}

//Handle bind address and added routing to this checker
func (hc *HealthChecker) Handle(addr, route string) error {
	hc.handler = NewHandler(hc, addr, route)
	return hc.handler.Server.ListenAndServe()
}

// NewHealthChecker return new instance HealthChecker with default parameters
func NewHealthChecker(name string) HealthChecker {
	return HealthChecker{
		Health: Health{
			Name:   name,
			Status: DOWN,
		},
		Checkers: make(map[string]Checker),
	}
}

// HealthError create new instance `Health` with `DOWN` status by error
func HealthError(err error) Health {
	return Health{
		Status: DOWN,
		Msg:    err.Error(),
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
	hc.SubHealth[name] = health
}
