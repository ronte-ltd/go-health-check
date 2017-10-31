package checkers

type FuncChecker struct {
	HealthChecker HealthChecker
	FuncCheck     func() Health
}

func SimpleChecker(name string) FuncChecker {
	return FuncChecker{
		HealthChecker: HealthChecker{
			Health: Health{
				Name:   name,
				Status: DOWN,
			},
			Checkers: make(map[string]Checker),
		},
	}
}

func NewFuncChecker(name string, check func() Health) FuncChecker {
	return FuncChecker{
		HealthChecker: HealthChecker{
			Health: Health{
				Name:   name,
				Status: DOWN,
			},
			Checkers: make(map[string]Checker),
		},
		FuncCheck: check,
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
		fc.HealthChecker.Status = UP
	}

	for _, c := range fc.HealthChecker.Checkers {
		var h, err = c.Check()
		if err != nil {
			fc.HealthChecker.SubHealth[c.Name()] = Health{
				Status: DOWN,
				Msg:    err.Error(),
			}
		} else {
			fc.HealthChecker.SubHealth[c.Name()] = h
			fc.checkStatus(h)
		}
	}
	if fc.FuncCheck != nil {
		var h = fc.FuncCheck()
		fc.HealthChecker.SubHealth[fc.HealthChecker.Name] = h
		fc.checkStatus(h)
	}
	return fc.HealthChecker.Health, nil
}

func (fc *FuncChecker) checkStatus(h Health) {
	if h.Status == DOWN {
		fc.HealthChecker.Status = DOWN
		fc.HealthChecker.Msg = h.Msg
	}
}

func (fc *FuncChecker) AddChecker(checker Checker) {
	fc.HealthChecker.Checkers[checker.Name()] = checker
}
