package scheduler

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

type mockMileageRunner struct {
	runs atomic.Int32
}

func (m *mockMileageRunner) RunMileageReminderAlerts(_ context.Context) (int, error) {
	m.runs.Add(1)
	return 1, nil
}

func TestParseMileageReminderInterval(t *testing.T) {
	cases := map[string]time.Duration{
		"":         DefaultMileageReminderInterval,
		"garbage":  DefaultMileageReminderInterval,
		"disabled": 0,
		"off":      0,
		"12h":      12 * time.Hour,
		"168h":     168 * time.Hour,
	}
	for raw, want := range cases {
		if got := ParseMileageReminderInterval(raw); got != want {
			t.Errorf("ParseMileageReminderInterval(%q) = %v, want %v", raw, got, want)
		}
	}
}

func TestMileageReminderCron_StartAndGracefulStop(t *testing.T) {
	runner := &mockMileageRunner{}
	cron := NewMileageReminderCron(runner, 15*time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cron.Start(ctx)
	time.Sleep(45 * time.Millisecond)
	cron.Stop()

	countAfterStop := runner.runs.Load()
	if countAfterStop < 1 {
		t.Fatalf("expected at least 1 scheduled execution before stop, got %d", countAfterStop)
	}

	time.Sleep(30 * time.Millisecond)
	if runner.runs.Load() != countAfterStop {
		t.Errorf("expected no executions after Stop(), before=%d after=%d", countAfterStop, runner.runs.Load())
	}
}

func TestMileageReminderCron_DisabledDoesNotStart(t *testing.T) {
	runner := &mockMileageRunner{}
	cron := NewMileageReminderCron(runner, 0)
	cron.Start(context.Background())
	time.Sleep(20 * time.Millisecond)
	cron.Stop()
	if runner.runs.Load() != 0 {
		t.Errorf("expected disabled cron not to run, got %d", runner.runs.Load())
	}
}
