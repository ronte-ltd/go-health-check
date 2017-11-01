package checkers

import (
	"fmt"
	"testing"
)

func TestHealthStatusUp(t *testing.T) {
	ch := NewHttpChecker("Yandex", "https://ya.ru")
	var health, err = ch.Check()
	if err != nil {
		t.Fatalf("Unhealthy, %s", err.Error())
	}
	if health.Status != UP {
		t.Fatalf("Yandex cannot be rechable, because %s", health.Msg)
	}
	t.Log(fmt.Sprintf("Healthy: %+v", health))
}

func TestCompositeHealthUp(t *testing.T) {
	yandex := NewHttpChecker("Yandex", "https://ya.ru")
	habr := NewHttpChecker("Habr", "https://habrahabr.ru")
	composite := CompositeChecker("Sites")
	composite.AddChecker(&yandex)
	composite.AddChecker(&habr)

	var health, err = composite.Check()
	if err != nil {
		t.Fatalf("Unhealthy, %s", err.Error())
	}
	if health.Status != UP {
		t.Fatalf("Yandex cannot be rechable, because %s", health.Msg)
	}
	t.Logf("Healthy: %+v", health)
}
