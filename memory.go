package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func readMemory() MemoryMetrics {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return MemoryMetrics{}
	}
	defer f.Close()

	var mem MemoryMetrics
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		key := strings.TrimSuffix(fields[0], ":")
		valKB, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			continue
		}
		valBytes := valKB * 1024

		switch key {
		case "MemTotal":
			mem.TotalBytes = valBytes
		case "MemAvailable":
			mem.AvailableBytes = valBytes
		case "MemFree":
			if mem.AvailableBytes == 0 {
				mem.AvailableBytes = valBytes
			}
		}
	}

	if mem.TotalBytes > 0 && mem.AvailableBytes <= mem.TotalBytes {
		mem.UsedBytes = mem.TotalBytes - mem.AvailableBytes
	}

	return mem
}


