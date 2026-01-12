package scheduler

import (
	"context"
	"time"

	"github.com/Govindpsd/api-monitoring-platform/internal/config"
	"github.com/Govindpsd/api-monitoring-platform/internal/probe"
)

func Start(
	ctx context.Context,
	p *probe.Probe,
	targets []config.Target,
	results chan<- probe.Result,
) {
	for _, target := range targets {
		go func(t config.Target) {
			ticker := time.NewTicker(t.Interval)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					result := p.Check(ctx, t.Name, t.URL)
					results <- result
				case <-ctx.Done():
					return
				}
			}
		}(target)
	}
}
