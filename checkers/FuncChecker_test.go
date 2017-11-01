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
	fc := NewFuncChecker("FuncChecker", check)
	var health, err = fc.Check()
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
	fc := NewFuncChecker("FuncChecker", check)
	var health, err = fc.Check()
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
	fc := NewFuncChecker("FuncChecker", check)
	fc1 := NewFuncChecker("FuncChecker1", check)
	fc2 := NewFuncChecker("FuncChecker2", check)
	fc.AddChecker(&fc1)
	fc.AddChecker(&fc2)
	var health, err = fc.Check()
	if err != nil {
		t.Fatalf("Should be healthy, %s", err.Error())
	}
	if health.Status != UP {
		t.Fatalf("Status should be UP")
	}
	t.Logf("Healthy: %+v", health)

}
