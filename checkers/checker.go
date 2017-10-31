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
