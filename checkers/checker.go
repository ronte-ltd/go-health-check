package checkers

import "github.com/go-kit/kit/log"

var (
	logger log.Logger
)

// `interface` for provide common method health checking
type Checker interface {
	Check() (Health, error)
	Name() string
}

type HealthChecker struct {
	Health
	Checkers map[string]Checker `json:"checkers,omitempty"`
}

type Health struct {
	Name      string            `json:"name"`
	Status    string            `json:"status"`
	Msg       string            `json:"msg,omitempty"`
	SubHealth map[string]Health `json:"subHealth,omitempty"`
}

const (
	DOWN = "DOWN"
	UP   = "UP"
)

// Create new HealthChecker with default parameters
func NewHealthChecker(name string) HealthChecker {
	return HealthChecker{
		Health: Health{
			Name:   name,
			Status: DOWN,
		},
		Checkers: make(map[string]Checker),
	}
}

// Create Health with `DOWN` status by error
func HealthError(err error) Health {
	return Health{
		Status: DOWN,
		Msg:    err.Error(),
	}
}

// Toggle Status to `UP`
func (hc *HealthChecker) Up() {
	hc.Status = UP
}

// Toggle Status to `DOWN`
func (hc *HealthChecker) Down() {
	hc.Status = DOWN
}

// Toggle Status to `DOWN` by error
func (hc *HealthChecker) DownError(err error) {
	hc.Status = DOWN
	hc.Msg = err.Error()
}

// Add SubHealth part for complex Health Checker
func (hc *HealthChecker) AddSubHealth(name string, health Health) {
	hc.SubHealth[name] = health
}
