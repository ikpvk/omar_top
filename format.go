package main

import "fmt"

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %siB", float64(bytes)/float64(div), []string{"", "K", "M", "G", "T", "P"}[exp])
}

func formatSpeed(bytesPerSec float64) string {
	if bytesPerSec < 0 {
		return "0 B/s"
	}
	const unit = 1024
	if bytesPerSec < unit {
		return fmt.Sprintf("%.0f B/s", bytesPerSec)
	}
	div, exp := float64(unit), 0
	for n := bytesPerSec / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %siB/s", bytesPerSec/div, []string{"", "K", "M", "G", "T", "P"}[exp])
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
