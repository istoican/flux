package flux

import (
	"sync"
	"testing"
)

func TestConcurrentMetrics(t *testing.T) {
	var metrics Metrics
	var wg sync.WaitGroup
	const numGoroutines = 1000

	// Increment in multiple goroutines
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			metrics.Keys.Increment()
		}()
	}
	wg.Wait()
	if metrics.Keys != numGoroutines {
		t.Errorf("expected %d, got %d", numGoroutines, metrics.Keys)
	}

	// Decrement in multiple goroutines
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			metrics.Keys.Decrement()
		}()
	}
	wg.Wait()
	if metrics.Keys != 0 {
		t.Errorf("expected 0, got %d", metrics.Keys)
	}
}

func TestMetricsKeysIncrement(t *testing.T) {
	var metrics Metrics

	metrics.Keys.Increment()
	if metrics.Keys != 1 {
		t.Errorf("expected 1, got %d", metrics.Keys)
	}
}

func TestMetricsInserts(t *testing.T) {
	var metrics Metrics

	metrics.Inserts.Increment()
	metrics.Inserts.Increment()
	if metrics.Inserts != 2 {
		t.Errorf("expected 2, got %d", metrics.Inserts)
	}
}

func TestMetricsDeletions(t *testing.T) {
	var metrics Metrics

	metrics.Deletions.Increment()
	metrics.Deletions.Decrement()
	if metrics.Deletions != 0 {
		t.Errorf("expected 0, got %d", metrics.Deletions)
	}
}

func TestMetricsReads(t *testing.T) {
	var metrics Metrics

	metrics.Reads.Set(5)
	if metrics.Reads != 5 {
		t.Errorf("expected 5, got %d", metrics.Reads)
	}
}
