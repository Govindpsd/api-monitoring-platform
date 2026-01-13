package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Govindpsd/api-monitoring-platform/internal/config"
	"github.com/Govindpsd/api-monitoring-platform/internal/metrics"
	"github.com/Govindpsd/api-monitoring-platform/internal/probe"
	"github.com/Govindpsd/api-monitoring-platform/internal/scheduler"
)

func main() {
	// 1️⃣ Global STOP button
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 2️⃣ Create Probe
	p := probe.NewProbe(5 * time.Second)

	// 3️⃣ Define targets (hardcoded approach)
	targets := []config.Target{
		{
			Name:     "google",
			URL:      "https://google.com",
			Interval: 5 * time.Second,
		},
		{
			Name:     "github",
			URL:      "https://github.com",
			Interval: 10 * time.Second,
		},
	}

	// 4️⃣ Results channel
	results := make(chan probe.Result)

	// 5️⃣ WaitGroup for scheduler workers
	var wg sync.WaitGroup

	// 5️⃣ Register metrics
	metrics.Register()

	// 6️⃣ Start scheduler
	go scheduler.Start(ctx, p, targets, results, &wg)

	// 7️⃣ Consume results
	go func() {
		for res := range results {
			metrics.ProbeTotal.WithLabelValues(res.Target).Inc()
			if res.Err != "" {
				metrics.ProbeFailures.WithLabelValues(res.Target).Inc()
				fmt.Printf(
					"❌ %s (%s) failed: %v\n",
					res.Target,
					res.URL,
					res.Err,
				)
				continue
			}
			// Record latency for successful probes
			metrics.ProbeLatency.WithLabelValues(res.Target).Observe(res.ResponseTime.Seconds())

			fmt.Printf(
				"✅ %s (%s) status=%d latency=%s\n",
				res.Target,
				res.URL,
				res.Status,
				res.ResponseTime,
			)
		}
	}()

	// 8️⃣ Health server
	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	server := &http.Server{
		Addr: ":8080",
	}

	go func() {
		fmt.Println("Health server running on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("HTTP server error:", err)
		}
	}()

	// 9️⃣ Wait for shutdown signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	fmt.Println("Shutdown signal received")

	// 10️⃣ Trigger graceful shutdown
	cancel()

	// Shutdown HTTP server
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	server.Shutdown(shutdownCtx)

	fmt.Println("Application shut down gracefully")
}
