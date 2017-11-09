package checkers

import (
	"testing"
)

func TestWithoutSubHealthUpStatus(t *testing.T) {
	var check = func() Health {
		return Health{
			Status: UP,
			Msg:    "All OK!",
		}
	}
	ch := NewHealthChecker("Main")
	ch.RegistryFunc("Up func", check)
	var health, err = ch.Check()
	if err != nil {
		t.Fatalf("Should be healthy, %s", err.Error())
	}
	if health.Status != UP {
		t.Errorf("Status should be UP")
	}
	t.Logf("Healthy: %+v", health)
}

func TestWithoutSubHealthDownStatus(t *testing.T) {
	var check = func() Health {
		return Health{
			Status: DOWN,
			Msg:    "Some bad...",
		}
	}
	hc := NewHealthChecker("FuncChecker")
	hc.RegistryFunc("Down", check)
	var health, err = hc.Check()
	if err != nil {
		t.Fatalf("Should be healthy, %s", err.Error())
	}
	if health.Status != DOWN {
		t.Fatalf("Status should be DOWN")
	}
	t.Logf("Healthy: %+v", health)
}

func TestTwoSubHealthUpStatus(t *testing.T) {
	var check = func() Health {
		return Health{
			Status: UP,
			Msg:    "All ok",
		}
	}
	fc := NewHealthChecker("Main")
	fc.RegistryFunc("FuncChecker1", check)
	fc.RegistryFunc("FuncChecker2", check)
	var health, err = fc.Check()
	if err != nil {
		t.Fatalf("Should be healthy, %s", err.Error())
	}
	if health.Status != UP {
		t.Fatalf("Status should be UP")
	}
	t.Logf("Healthy: %+v", health)

}
