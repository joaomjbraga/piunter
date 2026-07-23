package modules

import "testing"

func TestNewAppLogsModule(t *testing.T) {
	module := NewAppLogsModule()

	if module.ID() != "app-logs" {
		t.Fatalf("expected id 'app-logs', got %s", module.ID())
	}

	if module.Name() != "Logs de Aplicativos" {
		t.Fatalf("expected name 'Logs de Aplicativos', got %s", module.Name())
	}
}
