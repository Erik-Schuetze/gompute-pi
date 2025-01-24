package compute

import (
	"context"
	"testing"
	"time"
)

func TestComputePiTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	start := time.Now()
	result := ComputePi(ctx, int64(1000000))
	elapsed := time.Since(start)

	if result != nil {
		t.Errorf("expected nil result due to timeout, got %v", result)
	}

	if elapsed >= 2*time.Second {
		t.Errorf("expected timeout within 1 second, but took %v", elapsed)
	}
}
