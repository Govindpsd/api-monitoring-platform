package scheduler // Package scheduler handles periodic execution of HTTP probes

import (
	"context" // Context for cancellation and timeout management
	"sync"    // Sync provides WaitGroup for goroutine synchronization
	"time"    // Time provides ticker for periodic execution

	"github.com/Govindpsd/api-monitoring-platform/internal/config" // Config package contains Target struct
	"github.com/Govindpsd/api-monitoring-platform/internal/probe"  // Probe package contains Probe and Result types
)

func Start( // Start function initializes monitoring goroutines for each target
	ctx context.Context, // Context for cancellation signals
	p *probe.Probe, // Probe instance to perform HTTP checks
	targets []config.Target, // Slice of targets to monitor
	results chan<- probe.Result, // Channel to send probe results (send-only)
	wg *sync.WaitGroup, // WaitGroup to track goroutine completion
) {
	for _, target := range targets { // Iterate over each target in the slice
		wg.Add(1) // Increment WaitGroup counter for each goroutine
		// Start ONE goroutine per target
		go func(t config.Target) { // Launch a goroutine for each target (captures target by value)
			defer wg.Done() // Decrement WaitGroup counter when goroutine exits

			ticker := time.NewTicker(t.Interval) // Create a ticker that fires at the target's interval
			defer ticker.Stop()                  // Stop the ticker when goroutine exits to prevent leaks

			for { // Infinite loop to continuously monitor the target
				select { // Select statement to handle multiple channel operations
				case <-ticker.C: // Case: ticker fired, time to perform a check
					result := p.Check(ctx, t.Name, t.URL) // Execute HTTP probe check with target name and URL
					results <- result                     // Send the result to the results channel
				case <-ctx.Done(): // Case: context cancelled, shutdown signal received
					return // Exit the goroutine and stop monitoring this target
				}
			}
		}(target) // Pass target as argument to goroutine (captures by value to avoid closure issues)
	}
	// Close results channel when all workers exit
	go func() {
		wg.Wait()
		close(results)
	}()
}

/*
 * This scheduler package provides functionality to periodically monitor HTTP endpoints.
 *
 * The Start function:
 * - Takes a context for cancellation, a probe instance, a list of targets to monitor,
 *   a results channel, and a WaitGroup for synchronization
 * - Creates a separate goroutine for each target that will run independently
 * - Each goroutine uses a ticker to periodically execute HTTP probes at the target's interval
 * - Results are sent through the results channel for processing by the main application
 * - The WaitGroup ensures all goroutines are properly tracked and can be waited upon
 * - When the context is cancelled, all goroutines gracefully exit
 *
 * Key features:
 * - Concurrent monitoring of multiple targets with different intervals
 * - Graceful shutdown via context cancellation
 * - Resource cleanup with defer statements (ticker.Stop() and wg.Done())
 * - Closure-safe goroutine creation by passing target as a parameter
 */
