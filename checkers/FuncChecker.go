package checkers

type FuncChecker struct {
	HealthChecker HealthChecker
	FuncCheck     func() Health
}

func SimpleChecker(name string) FuncChecker {
	return FuncChecker{
		HealthChecker: NewHealthChecker(name),
	}
}

func NewFuncChecker(name string, check func() Health) FuncChecker {
	return FuncChecker{
		HealthChecker: NewHealthChecker(name),
		FuncCheck:     check,
	}
}

func (fc *FuncChecker) Name() string {
	return fc.HealthChecker.Name
}

func (fc *FuncChecker) Check() (Health, error) {
	if len(fc.HealthChecker.Checkers) == 0 {
		var health = fc.FuncCheck()
		health.Name = fc.HealthChecker.Name
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
	return fc.HealthChecker.Health, nil
}

func (fc *FuncChecker) selfCheck() {
	if fc.FuncCheck != nil {
		var h = fc.FuncCheck()
		fc.HealthChecker.AddSubHealth(fc.HealthChecker.Name, h)
		fc.checkDownStatus(h)
	}
}

func (fc *FuncChecker) checkDownStatus(h Health) {
	if h.Status == DOWN {
		fc.HealthChecker.Down()
		fc.HealthChecker.Msg = h.Msg
	}
}

func (fc *FuncChecker) AddChecker(checker Checker) {
	fc.HealthChecker.Checkers[checker.Name()] = checker
}
