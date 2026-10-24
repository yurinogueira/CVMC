package scheduler

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

type mockAlertRunner struct {
	runs atomic.Int32
}

func (m *mockAlertRunner) RunScheduledMaintenanceAlerts(_ context.Context) (int, error) {
	m.runs.Add(1)
	return 1, nil
}

func TestParseCronInterval(t *testing.T) {
	if got := ParseCronInterval(""); got != DefaultCronInterval {
		t.Errorf("expected default %v, got %v", DefaultCronInterval, got)
	}
	if got := ParseCronInterval("disabled"); got != 0 {
		t.Errorf("expected 0 for disabled, got %v", got)
	}
	if got := ParseCronInterval("6h"); got != 6*time.Hour {
		t.Errorf("expected 6h, got %v", got)
	}
}

func TestMaintenanceCron_StartAndGracefulStop(t *testing.T) {
	runner := &mockAlertRunner{}
	cron := NewMaintenanceCron(runner, 15*time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cron.Start(ctx)
	time.Sleep(45 * time.Millisecond)
	cron.Stop()

	countAfterStop := runner.runs.Load()
	if countAfterStop < 1 {
		t.Fatalf("expected at least 1 scheduled execution before stop, got %d", countAfterStop)
	}

	// Verify no further runs happen after Stop()
	time.Sleep(30 * time.Millisecond)
	if runner.runs.Load() != countAfterStop {
		t.Errorf("expected no executions after Stop(), before=%d after=%d", countAfterStop, runner.runs.Load())
	}
}
