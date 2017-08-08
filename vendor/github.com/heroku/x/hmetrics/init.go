package hmetrics

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	metricWaitTime = 20 * time.Second
)

type printfer interface {
	Printf(s string, v ...interface{})
}

type AlreadyStarted struct{}

func (as AlreadyStarted) Error() string {
	return "already started"
}

func (as AlreadyStarted) Fatal() bool {
	return false
}

type HerokuMetricsURLUnset struct{}

func (e HerokuMetricsURLUnset) Error() string {
	return "cannot report metrics because HEROKU_METRICS_URL is unset"
}

func (e HerokuMetricsURLUnset) Fatal() bool {
	return true
}

var (
	mu      sync.Mutex
	started bool
)

func Report(ctx context.Context, pfer printfer) error {
	mu.Lock()
	defer mu.Unlock()
	if started {
		return AlreadyStarted{}
	}
	endpoint := os.Getenv("HEROKU_METRICS_URL")
	if endpoint == "" {
		return HerokuMetricsURLUnset{}
	}
	go report(ctx, endpoint, pfer)
	started = true
	return nil
}

func report(ctx context.Context, endpoint string, pfer printfer) {
	t := time.NewTicker(metricWaitTime)
	defer t.Stop()

	for {
		select {
		case <-t.C:
		case <-ctx.Done():
			mu.Lock()
			defer mu.Unlock()
			started = false
			return
		}

		if err := gatherMetrics(); err != nil && pfer != nil {
			pfer.Printf("error encoding metrics: %v", err)
			continue
		}
		if err := submitPayload(ctx, endpoint); err != nil && pfer != nil {
			pfer.Printf("error submitting metrics: %v", err)
			continue
		}
	}
}

var (
	lastGCPause uint64
	buf         bytes.Buffer
)

func gatherMetrics() error {
	stats := &runtime.MemStats{}
	runtime.ReadMemStats(stats)

	// cribbed from https://github.com/codahale/metrics/blob/master/runtime/memstats.go

	pauseNS := stats.PauseTotalNs - lastGCPause
	lastGCPause = stats.PauseTotalNs

	result := struct {
		Counters map[string]float64 `json:"counters"`
		Gauges   map[string]float64 `json:"gauges"`
	}{
		Counters: map[string]float64{
			"go.gc.collections": float64(stats.NumGC),
			"go.gc.pause.ns":    float64(pauseNS),
		},
		Gauges: map[string]float64{
			"go.memory.heap.bytes":  float64(stats.Alloc),
			"go.memory.stack.bytes": float64(stats.StackInuse),
		},
	}

	buf.Reset()
	return json.NewEncoder(&buf).Encode(result)
}

func submitPayload(ctx context.Context, where string) error {
	req, err := http.NewRequest("POST", where, &buf)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected %v (http.StatusOK) but got %s", http.StatusOK, resp.Status)
	}

	return nil
}
