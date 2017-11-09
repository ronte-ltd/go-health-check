package checkers

import (
	"context"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	first := MockHTTPServer("22922")
	second := MockHTTPServer("22923")
	resp := m.Run()
	first.Shutdown(context.Background())
	second.Shutdown(context.Background())
	os.Exit(resp)
}

func TestHealthStatusUp(t *testing.T) {
	ch := NewHealthChecker("Main")
	ch.RegistryURL("First", "http://127.0.0.1:22922/22922")
	var health, err = ch.Check()
	if err != nil {
		t.Fatalf("Unhealthy, %s", err.Error())
	}
	if health.Status != UP {
		t.Fatalf("cannot be rechable, because %s", health.Msg)
	}
	t.Logf("Healthy: %+v", health)
}

func TestCompositeHealthUp(t *testing.T) {
	ch := NewHealthChecker("Main")
	ch.RegistryURL("Yandex", "http://127.0.0.1:22922/22922")
	ch.RegistryURL("Habr", "http://127.0.0.1:22923/22923")

	var health, err = ch.Check()
	if err != nil {
		t.Fatalf("Unhealthy, %s", err.Error())
	}
	if health.Status != UP {
		t.Fatalf("Yandex cannot be rechable, because %s", health.Msg)
	}
	t.Logf("Healthy: %+v", health)
}
