package probe

import (
	"context"
	"net/http"
	"time"
)

// result contains the result of a http probe
type Result struct {
	URL          string
	Target       string
	Status       int
	ResponseTime time.Duration
	Err          string
}
type Probe struct { //it holds the configuration for a http probe
	client *http.Client //http client to make the http request

}

func NewProbe(timeout time.Duration) *Probe { //create a new probe with reusable http client
	return &Probe{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}
func (p *Probe) Check(
	ctx context.Context,
	url string,
	target string,

) Result {
	start := time.Now()
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return Result{
			URL:          url,
			Status:       0,
			ResponseTime: 0,
			Err:          err.Error(),
		}
	}
	//send request using the same client to avoid creating a new one for each request
	resp, err := p.client.Do(req)
	if err != nil {
		return Result{
			URL:    url,
			Target: target,
			Err:    err.Error(),
		}
	}
	defer resp.Body.Close()
	return Result{
		URL:          url,
		Target:       target,
		Status:       resp.StatusCode,
		ResponseTime: time.Since(start),
	}
}
