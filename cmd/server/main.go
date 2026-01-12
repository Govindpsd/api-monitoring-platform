package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Govindpsd/api-monitoring-platform/internal/config"
	"github.com/Govindpsd/api-monitoring-platform/internal/probe"
	"github.com/Govindpsd/api-monitoring-platform/internal/scheduler"
)

func main() {
	targets := []config.Target{ //create a new config with the targets
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
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	//root context for the application
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel the context when the application is shutdown

	//create a new probe
	p := probe.NewProbe(5 * time.Second)

	//Create a channel to receive the results of the probes
	results := make(chan probe.Result)

	//start the scheduler
	go scheduler.Start(ctx, p, targets, results)

	//listen for shutdown signals
	go func() {
		for res := range results {
			if res.Err != "nil" {
				fmt.Println("Error: ", res.Target, "Failed with error: ", res.Err)
			} else {
				fmt.Printf(
					"❌ %s (%s) failed with error: %v\n",
					res.Target,
					res.URL,
					res.Err,
				)
			}
		}
		go func() {
			fmt.Println("Health server running on :8080")
			http.ListenAndServe(":8080", nil)
		}()
	}()

	//wait for the context to be cancelled
	sigCh := make(chan os.Signal, 1)                      // channel to listen for shutdown signals
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM) // listen for interrupt and term signals
	<-sigCh                                               // wait for a signal to shutdown the application
	println("shut down signal received")
	cancel() // cancel the context
	fmt.Println("Application shut down")
}
