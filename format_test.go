package main

import (
	"math"
	"testing"
)

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		input uint64
		want  string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.0 KiB"},
		{1536, "1.5 KiB"},
		{1048576, "1.0 MiB"},
		{1073741824, "1.0 GiB"},
		{1099511627776, "1.0 TiB"},
		{1125899906842624, "1.0 PiB"},
	}
	for _, tt := range tests {
		got := formatBytes(tt.input)
		if got != tt.want {
			t.Errorf("formatBytes(%d) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFormatBytesMax(t *testing.T) {
	// Max uint64 should not panic (exp capped at len(suffixes)-1)
	got := formatBytes(math.MaxUint64)
	if got == "" {
		t.Error("formatBytes(MaxUint64) returned empty string")
	}
}

func TestFormatSpeed(t *testing.T) {
	tests := []struct {
		input float64
		want  string
	}{
		{-1, "0 B/s"},
		{0, "0 B/s"},
		{500, "500 B/s"},
		{1000, "1.0 kB/s"},
		{1500, "1.5 kB/s"},
		{1000000, "1.0 MB/s"},
		{1000000000, "1.0 GB/s"},
		{1000000000000, "1.0 TB/s"},
		{1000000000000000, "1.0 PB/s"},
	}
	for _, tt := range tests {
		got := formatSpeed(tt.input)
		if got != tt.want {
			t.Errorf("formatSpeed(%v) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFormatSpeedMax(t *testing.T) {
	// Very large value should not panic (exp capped at len(suffixes)-1)
	got := formatSpeed(1e19)
	if got == "" {
		t.Error("formatSpeed(1e19) returned empty string")
	}
}

func TestFormatFreq(t *testing.T) {
	tests := []struct {
		input float64
		want  string
	}{
		{0, "0 MHz"},
		{999, "999 MHz"},
		{1000, "1.0 GHz"},
		{2500, "2.5 GHz"},
	}
	for _, tt := range tests {
		got := formatFreq(tt.input)
		if got != tt.want {
			t.Errorf("formatFreq(%v) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFormatTemp(t *testing.T) {
	got := formatTemp(75.3)
	want := "75°C"
	if got != want {
		t.Errorf("formatTemp(75.3) = %q, want %q", got, want)
	}
}

func TestFormatPct(t *testing.T) {
	tests := []struct {
		input float64
		want  string
	}{
		{0, "0%"},
		{50, "50%"},
		{100, "100%"},
	}
	for _, tt := range tests {
		got := formatPct(tt.input)
		if got != tt.want {
			t.Errorf("formatPct(%v) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
