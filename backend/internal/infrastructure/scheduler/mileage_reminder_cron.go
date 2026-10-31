package scheduler

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"
)

// DefaultMileageReminderInterval is the default (weekly) interval for the mileage reminder cron.
const DefaultMileageReminderInterval = 7 * 24 * time.Hour

type MileageReminderRunner interface {
	RunMileageReminderAlerts(ctx context.Context) (int, error)
}

// MileageReminderCron periodically reminds users to update the mileage of stale vehicles.
type MileageReminderCron struct {
	runner   MileageReminderRunner
	interval time.Duration
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	mu       sync.Mutex
	running  bool
}

// ParseMileageReminderInterval parses MILEAGE_REMINDER_CRON_INTERVAL, defaulting to a weekly interval.
func ParseMileageReminderInterval(raw string) time.Duration {
	cleaned := strings.ToLower(strings.TrimSpace(raw))
	switch cleaned {
	case "":
		return DefaultMileageReminderInterval
	case "0", "off", "false", "disabled", "none":
		return 0
	}
	if d, err := time.ParseDuration(cleaned); err == nil && d > 0 {
		return d
	}
	return DefaultMileageReminderInterval
}

func NewMileageReminderCron(runner MileageReminderRunner, interval time.Duration) *MileageReminderCron {
	return &MileageReminderCron{
		runner:   runner,
		interval: interval,
	}
}

func (c *MileageReminderCron) Start(parentCtx context.Context) {
	if c == nil || c.runner == nil || c.interval <= 0 {
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

func (c *MileageReminderCron) RunOnce(ctx context.Context) (int, error) {
	if c == nil || c.runner == nil {
		return 0, nil
	}
	runCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	return c.runner.RunMileageReminderAlerts(runCtx)
}

func (c *MileageReminderCron) runOnce(ctx context.Context) {
	sent, err := c.RunOnce(ctx)
	if err != nil {
		log.Printf("[mileage-reminder-cron] error evaluating mileage reminders: %v", err)
		return
	}
	if sent > 0 {
		log.Printf("[mileage-reminder-cron] dispatched %d mileage reminder(s)", sent)
	}
}

func (c *MileageReminderCron) Stop() {
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
