package scheduler

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"
)

const DefaultCronInterval = 24 * time.Hour

type AlertRunner interface {
	RunScheduledMaintenanceAlerts(ctx context.Context) (int, error)
}

type MaintenanceCron struct {
	runner   AlertRunner
	interval time.Duration
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	mu       sync.Mutex
	running  bool
}

func ParseCronInterval(raw string) time.Duration {
	cleaned := strings.ToLower(strings.TrimSpace(raw))
	if cleaned == "" {
		return DefaultCronInterval
	}
	switch cleaned {
	case "0", "off", "false", "disabled", "none":
		return 0
	}
	if d, err := time.ParseDuration(cleaned); err == nil && d > 0 {
		return d
	}
	return DefaultCronInterval
}

func NewMaintenanceCron(runner AlertRunner, interval time.Duration) *MaintenanceCron {
	return &MaintenanceCron{
		runner:   runner,
		interval: interval,
	}
}

func (c *MaintenanceCron) Start(parentCtx context.Context) {
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

func (c *MaintenanceCron) RunOnce(ctx context.Context) (int, error) {
	if c == nil || c.runner == nil {
		return 0, nil
	}
	runCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	return c.runner.RunScheduledMaintenanceAlerts(runCtx)
}

func (c *MaintenanceCron) runOnce(ctx context.Context) {
	alerts, err := c.RunOnce(ctx)
	if err != nil {
		log.Printf("[maintenance-cron] error evaluating preventive alerts: %v", err)
		return
	}
	if alerts > 0 {
		log.Printf("[maintenance-cron] dispatched %d consolidated maintenance alert(s)", alerts)
	}
}

func (c *MaintenanceCron) Stop() {
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
