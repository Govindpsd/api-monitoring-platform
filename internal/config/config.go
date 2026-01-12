package config

import "time"

// Target contains the on API to monitor
type Target struct {
	URL      string
	Name     string
	Target   string
	Interval time.Duration
}

// Config contains the configuration for the application
type Config struct {
	Targets []Target
}
