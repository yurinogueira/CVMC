package scheduler

import (
	"context"
	"sync"
	"testing"
	"time"
)

type mockCleaner struct {
	mu     sync.Mutex
	called int
}

func (m *mockCleaner) CleanupExpiredReports(ctx context.Context, retentionDuration time.Duration, batchSize int) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.called++
	return 1, nil
}

func TestReportCleanupCron_StartAndStop(t *testing.T) {
	cleaner := &mockCleaner{}
	cron := NewReportCleanupCron(cleaner, 15*time.Millisecond, 30*24*time.Hour)

	ctx := context.Background()
	cron.Start(ctx)

	time.Sleep(50 * time.Millisecond)
	cron.Stop()

	cleaner.mu.Lock()
	called := cleaner.called
	cleaner.mu.Unlock()

	if called == 0 {
		t.Fatalf("expected cleaner to be called, got %d", called)
	}

	purged, err := cron.RunOnce(ctx)
	if err != nil || purged != 1 {
		t.Fatalf("unexpected RunOnce result: %d, err: %v", purged, err)
	}
}
