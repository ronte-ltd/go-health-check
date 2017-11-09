// Package checkers provide health check for different service
package checkers

// FuncChecker provide health check via `func()`
type FuncChecker struct {
	HealthChecker HealthChecker
	FuncCheck     func() Health
}

// NewCompositeChecker Create Composite Checker without self-check
func NewCompositeChecker(name string) FuncChecker {
	return FuncChecker{
		HealthChecker: NewHealthChecker(name),
	}
}

// NewFuncChecker return Composite Checker with self-check func
func NewFuncChecker(name string, check func() Health) *FuncChecker {
	return &FuncChecker{
		HealthChecker: NewHealthChecker(name),
		FuncCheck:     check,
	}
}

//Name return name current Health check
func (fc *FuncChecker) Name() string {
	return fc.HealthChecker.Name()
}

//Check execute health check func and sub-func and return Health or error
func (fc *FuncChecker) Check() (Health, error) {
	if len(fc.HealthChecker.Checkers) == 0 {
		health := fc.FuncCheck()
		health.Name = fc.HealthChecker.Name()
		fc.HealthChecker.PushHealth()
		return health, nil
	}

	if fc.HealthChecker.SubHealth == nil {
		fc.HealthChecker.SubHealth = make(map[string]Health, len(fc.HealthChecker.Checkers)+1)
		fc.HealthChecker.Up()
	}

	for _, c := range fc.HealthChecker.Checkers {
		var h, err = c.Check()
		if err != nil {
			h = HealthError(err)
		}
		fc.HealthChecker.AddSubHealth(c.Name(), h)
		fc.checkDownStatus(h)
	}

	fc.selfCheck()
	fc.HealthChecker.PushHealth()
	return fc.HealthChecker.Health, nil
}

func (fc *FuncChecker) selfCheck() {
	if fc.FuncCheck != nil {
		h := fc.FuncCheck()
		fc.HealthChecker.AddSubHealth(fc.HealthChecker.Name(), h)
		fc.checkDownStatus(h)
	}
}

func (fc *FuncChecker) checkDownStatus(h Health) {
	if h.Status == DOWN {
		fc.HealthChecker.Down()
		fc.HealthChecker.Msg = h.Msg
	}
}

// AddChecker add new Sub-checker to Composite Checker
func (fc *FuncChecker) AddChecker(checker Checker) {
	fc.HealthChecker.Checkers[checker.Name()] = checker
}
