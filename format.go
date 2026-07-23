package main

import "fmt"

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	suffixes := []string{"", "K", "M", "G", "T", "P"}
	div := uint64(1)
	exp := 0
	for exp < len(suffixes)-1 && bytes/div >= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %siB", float64(bytes)/float64(div), suffixes[exp])
}

func formatSpeed(bytesPerSec float64) string {
	if bytesPerSec < 0 {
		return "0 B/s"
	}
	const unit = 1000
	if bytesPerSec < unit {
		return fmt.Sprintf("%.0f B/s", bytesPerSec)
	}
	suffixes := []string{"", "k", "M", "G", "T", "P"}
	div := float64(1)
	exp := 0
	for exp < len(suffixes)-1 && bytesPerSec/div >= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %sB/s", bytesPerSec/div, suffixes[exp])
}

func formatFreq(mhz float64) string {
	if mhz < 1000 {
		return fmt.Sprintf("%.0f MHz", mhz)
	}
	return fmt.Sprintf("%.1f GHz", mhz/1000)
}

func formatTemp(celsius float64) string {
	return fmt.Sprintf("%.0f°C", celsius)
}

func formatPct(pct float64) string {
	return fmt.Sprintf("%.0f%%", pct)
}
