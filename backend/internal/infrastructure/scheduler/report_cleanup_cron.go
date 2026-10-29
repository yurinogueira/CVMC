package scheduler

import (
	"context"
	"log"
	"sync"
	"time"
)

type ReportCleaner interface {
	CleanupExpiredReports(ctx context.Context, retentionDuration time.Duration, batchSize int) (int, error)
}

type ReportCleanupCron struct {
	cleaner   ReportCleaner
	interval  time.Duration
	retention time.Duration
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	mu        sync.Mutex
	running   bool
}

func NewReportCleanupCron(cleaner ReportCleaner, interval time.Duration, retention time.Duration) *ReportCleanupCron {
	if retention <= 0 {
		retention = 30 * 24 * time.Hour
	}
	return &ReportCleanupCron{
		cleaner:   cleaner,
		interval:  interval,
		retention: retention,
	}
}

func (c *ReportCleanupCron) Start(parentCtx context.Context) {
	if c == nil || c.cleaner == nil || c.interval <= 0 {
		return
	}

	c.mu.Lock()
	if c.running {
		c.mu.Unlock()
		return
	}
	ctx, cancel := context.WithCancel(parentCtx)
	c.cancel = cancel
	c.running = true
	c.wg.Add(1)
	c.mu.Unlock()

	go func() {
		defer c.wg.Done()
		ticker := time.NewTicker(c.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				c.runOnce(ctx)
			}
		}
	}()
}

func (c *ReportCleanupCron) RunOnce(ctx context.Context) (int, error) {
	if c == nil || c.cleaner == nil {
		return 0, nil
	}
	runCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	return c.cleaner.CleanupExpiredReports(runCtx, c.retention, 50)
}

func (c *ReportCleanupCron) runOnce(ctx context.Context) {
	purged, err := c.RunOnce(ctx)
	if err != nil {
		log.Printf("[report-cleanup-cron] error running cleanup: %v", err)
		return
	}
	if purged > 0 {
		log.Printf("[report-cleanup-cron] successfully purged %d expired report(s)", purged)
	}
}

func (c *ReportCleanupCron) Stop() {
	if c == nil {
		return
	}
	c.mu.Lock()
	cancel := c.cancel
	wasRunning := c.running
	c.running = false
	c.cancel = nil
	c.mu.Unlock()

	if cancel != nil {
		cancel()
	}
	if wasRunning {
		c.wg.Wait()
	}
}
