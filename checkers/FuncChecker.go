// Package checkers provide health check for different service
package checkers

import "sync"

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
	if fc.HealthChecker.Checkers.Len() == 0 {
		health := fc.FuncCheck()
		health.Name = fc.HealthChecker.Name()
		fc.HealthChecker.PushHealth()
		return health, nil
	}

	if fc.HealthChecker.SubHealth == nil {
		fc.HealthChecker.SubHealth = NewSubHealthMapWithLen(fc.HealthChecker.Checkers.Len() + 1)
		fc.HealthChecker.Up()
	}

	wg := sync.WaitGroup{}
	wg.Add(fc.HealthChecker.Checkers.Len())

	f := func(key string, value Checker) bool {

		var h, err = value.Check()
		if err != nil {
			h = HealthError(err)
		}
		fc.HealthChecker.AddSubHealth(value.Name(), h)
		fc.checkDownStatus(h)
		wg.Done()
		return true
	}
	go fc.HealthChecker.Checkers.Range(f)

	wg.Wait()
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
	fc.HealthChecker.Checkers.Store(checker.Name(), checker)
}
