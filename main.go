package main

import (
	"context"
	"erik-schuetze/gompute-pi/compute"
	"flag"
	"fmt"
	"math/big"
	"time"
)

type Config struct {
	Iterations int64
	MaxTime    int64
	Precision  int
}

type Result struct {
	Value     *big.Float
	Precision int
	Error     error
}

func calculatePi(ctx context.Context, cfg Config) Result {
	fmt.Printf("[%s] Starting calculation\n", time.Now().Format(time.RFC3339))
	resultChan := make(chan Result, 1)

	go func() {
		fmt.Printf("[%s] Starting goroutine\n", time.Now().Format(time.RFC3339))
		value := compute.ComputePi(ctx, cfg.Iterations)
		if value == nil {
			resultChan <- Result{Error: ctx.Err()}
			return
		}
		resultChan <- Result{
			Value:     value,
			Precision: cfg.Precision,
		}
	}()

	select {
	case result := <-resultChan:
		fmt.Printf("[%s] Received result\n%v\n", time.Now().Format(time.RFC3339), result.Value.Text('f', 50))
		return result
	case <-ctx.Done():
		fmt.Printf("[%s] Context cancelled\n", time.Now().Format(time.RFC3339))
		return Result{Error: ctx.Err()}
	}
}

func setupContext(maxTime int64) (context.Context, context.CancelFunc) {
	if maxTime <= 0 {
		fmt.Println("Setting up context without timeout")
		return context.WithCancel(context.Background())
	}
	fmt.Println("Setting up context with timeout")
	return context.WithTimeout(context.Background(), time.Duration(maxTime)*time.Second)
}

func run(cfg Config) (string, error) {
	ctx, cancel := setupContext(cfg.MaxTime)
	defer cancel()

	result := calculatePi(ctx, cfg)
	if result.Error != nil {
		return "", result.Error
	}

	return result.Value.Text('f', result.Precision), nil
}

func main() {
	// define flags for iterations, maxTime
	iterationsPtr := flag.Int("i", 1, "number of iterations to calculate Pi")
	maxTimePtr := flag.Int("t", 0, "maximum time to calculate Pi in seconds")
	precisionPtr := flag.Int("p", 50, "presicion of the printed result")

	flag.Parse()

	cfg := Config{
		Iterations: int64(*iterationsPtr),
		MaxTime:    int64(*maxTimePtr),
		Precision:  int(*precisionPtr),
	}

	result, err := run(cfg)
	if err != nil {
		fmt.Printf("Calculation failed: %v\n", err)
		return
	}
	fmt.Printf("Calculation completed: %s\n", result)
}
