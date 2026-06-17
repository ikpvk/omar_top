package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func readCPUStat() cpuData {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return cpuData{}
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "cpu ") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 5 {
			return cpuData{}
		}
		var total uint64
		for i := 1; i < len(fields); i++ {
			val, _ := strconv.ParseUint(fields[i], 10, 64)
			total += val
		}
		idle, _ := strconv.ParseUint(fields[4], 10, 64)
		return cpuData{total: total, idle: idle}
	}
	return cpuData{}
}

func readCPUCores() int {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return 0
	}
	defer f.Close()
	var cores int
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu") {
			fields := strings.Fields(line)
			if len(fields) > 0 {
				name := fields[0]
				if name != "cpu" {
					_, err := fmt.Sscanf(name, "cpu%d", &cores)
					if err == nil {
						cores++
					}
				}
			}
		}
	}
	return cores
}

func readCPUFreq() float64 {
	f, err := os.ReadFile("/sys/devices/system/cpu/cpu0/cpufreq/scaling_cur_freq")
	if err != nil {
		f, err = os.ReadFile("/proc/cpuinfo")
		if err != nil {
			return 0
		}
		scanner := bufio.NewScanner(strings.NewReader(string(f)))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "cpu MHz") {
				parts := strings.Split(line, ":")
				if len(parts) == 2 {
					mhz, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
					return mhz
				}
			}
		}
		return 0
	}
	mhz, _ := strconv.ParseFloat(strings.TrimSpace(string(f)), 64)
	return mhz
}

func readCPUTemp() float64 {
	base := "/sys/class/thermal/thermal_zone0/temp"
	data, err := os.ReadFile(base)
	if err != nil {
		return 0
	}
	millicelsius, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
	if err != nil {
		return 0
	}
	return math.Round(millicelsius/10) / 100
}
