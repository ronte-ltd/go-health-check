package checkers

type Checker interface {
	Check() (Health, error)
	Name() string
}

type HealthChecker struct {
	Health   `json:"health"`
	Checkers map[string]Checker `json:"checkers"`
}

type Health struct {
	Name      string            `json:"name"`
	Status    string            `json:"status"`
	Msg       string            `json:"msg, omitempty"`
	SubHealth map[string]Health `json:"health, omitempty"`
}

const (
	DOWN = "DOWN"
	UP   = "UP"
)

func NewHealthChecker(name string) HealthChecker {
	return HealthChecker{
		Health: Health{
			Name:   name,
			Status: DOWN,
		},
		Checkers: make(map[string]Checker),
	}
}

func (hc *HealthChecker) Up() {
	hc.Status = UP
}

func (hc *HealthChecker) Down() {
	hc.Status = DOWN
}
