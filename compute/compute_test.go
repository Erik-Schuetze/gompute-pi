package compute

import (
	"context"
	"math"
	"testing"
)

func TestComputePiOneIteration(t *testing.T) {
	ctx := context.Background()
	result := ComputePi(ctx, 1)
	pi, _ := result.Float64()
	if math.Abs(pi-math.Pi) > 0.0001 {
		t.Errorf("Expected Pi to be close to %v, got %v", math.Pi, pi)
	}
}

func TestComputePiTenIteration(t *testing.T) {
	ctx := context.Background()
	result := ComputePi(ctx, 10)
	pi, _ := result.Float64()
	if math.Abs(pi-math.Pi) > 0.0000001 {
		t.Errorf("Expected Pi to be close to %v, got %v", math.Pi, pi)
	}
}

func TestComputePiHundredIterations(t *testing.T) {
	ctx := context.Background()
	result := ComputePi(ctx, 100)
	pi, _ := result.Float64()
	if math.Abs(pi-math.Pi) > 1e-8 {
		t.Errorf("Expected Pi to be close to %v, got %v", math.Pi, pi)
	}
}
