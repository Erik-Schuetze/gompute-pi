package main

import (
	"testing"
)

func TestCalculatePi(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		want    float64
		wantErr bool
	}{
		{
			name: "basic calculation",
			cfg: Config{
				Iterations: 1,
				MaxTime:    5,
				Precision:  10,
			},
			wantErr: false,
		},
		{
			name: "timeout",
			cfg: Config{
				Iterations: 1000000,
				MaxTime:    1,
				Precision:  10,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := run(tt.cfg)

			if (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result == "" {
				t.Error("expected result, got empty string")
			}
		})
	}
}

func TestSetupContext(t *testing.T) {
	tests := []struct {
		name    string
		maxTime int64
		want    bool
	}{
		{"no timeout", 0, false},
		{"with timeout", 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := setupContext(tt.maxTime)
			defer cancel()

			_, hasDeadline := ctx.Deadline()
			if hasDeadline != tt.want {
				t.Errorf("setupContext() deadline = %v, want %v", hasDeadline, tt.want)
			}
		})
	}
}

func TestTimeOut(t *testing.T) {
	cfg := Config{
		Iterations: 10000,
		MaxTime:    1,
		Precision:  10,
	}

	result, err := run(cfg)

	if err == nil {
		t.Error("expected timeout error, got nil")
	}

	if result != "" {
		t.Errorf("expected empty result due to timeout, got %v", result)
	}
}
