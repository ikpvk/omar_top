package main

import (
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type gpuVendor int

const (
	gpuVendorNone gpuVendor = iota
	gpuVendorNVIDIA
	gpuVendorAMD
	gpuVendorIntel
)

func detectGPU() (gpuVendor, string) {
	cardDir := "/sys/class/drm"
	entries, err := os.ReadDir(cardDir)
	if err != nil {
		return gpuVendorNone, ""
	}
	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasPrefix(name, "card") {
			continue
		}
		vendorPath := filepath.Join(cardDir, name, "device", "vendor")
		data, err := os.ReadFile(vendorPath)
		if err != nil {
			continue
		}
		vendor := strings.TrimSpace(string(data))
		switch vendor {
		case "0x10de":
			return gpuVendorNVIDIA, "NVIDIA"
		case "0x1002":
			return gpuVendorAMD, "AMD"
		case "0x8086":
			return gpuVendorIntel, "Intel"
		}
	}
	return gpuVendorNone, ""
}

func getGPU() *GPUMetrics {
	vendor, name := detectGPU()
	if vendor == gpuVendorNone {
		return nil
	}
	switch vendor {
	case gpuVendorNVIDIA:
		return getNvidiaGPU(name)
	case gpuVendorAMD:
		return getAMDGpu(name)
	case gpuVendorIntel:
		return getIntelGpu(name)
	}
	return nil
}

func getNvidiaGPU(name string) *GPUMetrics {
	cmd := exec.Command("nvidia-smi",
		"--query-gpu=utilization.gpu,temperature.gpu,memory.used,memory.total",
		"--format=csv,noheader,nounits")
	output, err := cmd.Output()
	if err != nil {
		return nil
	}
	parts := strings.Split(strings.TrimSpace(string(output)), ", ")
	if len(parts) < 2 {
		return nil
	}
	usage, _ := strconv.ParseFloat(parts[0], 64)
	temp, _ := strconv.ParseFloat(parts[1], 64)

	gpu := &GPUMetrics{
		UsagePercent: usage,
		Temperature:  temp,
		Name:         name,
	}

	if len(parts) >= 4 {
		memUsed, _ := strconv.ParseUint(parts[2], 10, 64)
		memTotal, _ := strconv.ParseUint(parts[3], 10, 64)
		gpu.HasMemory = true
		gpu.MemUsed = memUsed * 1024 * 1024
		gpu.MemTotal = memTotal * 1024 * 1024
	}

	return gpu
}

func getAMDGpu(name string) *GPUMetrics {
	usage := readGPUUtilization("amdgpu")
	if usage < 0 {
		return nil
	}
	temp := readGPUTemp("amdgpu")
	return &GPUMetrics{
		UsagePercent: usage,
		Temperature:  temp,
		Name:         name,
	}
}

func getIntelGpu(name string) *GPUMetrics {
	usage := readGPUUtilization("i915")
	if usage < 0 {
		return nil
	}
	temp := readGPUTemp("i915")
	return &GPUMetrics{
		UsagePercent: usage,
		Temperature:  temp,
		Name:         name,
	}
}

func readGPUUtilization(driver string) float64 {
	data, err := os.ReadFile("/sys/class/drm/card0/device/gpu_busy_percent")
	if err != nil {
		entries, _ := filepath.Glob("/sys/class/drm/card*/device/gpu_busy_percent")
		if len(entries) > 0 {
			data, err = os.ReadFile(entries[0])
		}
	}
	if err != nil {
		return -1
	}
	pct, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
	if err != nil {
		return -1
	}
	return pct
}

func readGPUTemp(driver string) float64 {
	hwmonBase := "/sys/class/hwmon"
	entries, err := os.ReadDir(hwmonBase)
	if err != nil {
		return 0
	}
	for _, entry := range entries {
		namePath := filepath.Join(hwmonBase, entry.Name(), "name")
		data, err := os.ReadFile(namePath)
		if err != nil {
			continue
		}
		if strings.TrimSpace(string(data)) == driver {
			tempPath := filepath.Join(hwmonBase, entry.Name(), "temp1_input")
			tempData, err := os.ReadFile(tempPath)
			if err != nil {
				continue
			}
			milli, _ := strconv.ParseFloat(strings.TrimSpace(string(tempData)), 64)
			return math.Round(milli/10) / 100
		}
	}

	entries, err = os.ReadDir(hwmonBase)
	if err != nil {
		return 0
	}
	for _, entry := range entries {
		labelPath := filepath.Join(hwmonBase, entry.Name(), "temp1_label")
		data, err := os.ReadFile(labelPath)
		if err != nil {
			continue
		}
		label := strings.TrimSpace(string(data))
		if strings.Contains(strings.ToLower(label), "gpu") || strings.Contains(strings.ToLower(label), "edge") {
			tempInput := filepath.Join(hwmonBase, entry.Name(), "temp1_input")
			tempData, err := os.ReadFile(tempInput)
			if err != nil {
				continue
			}
			milli, _ := strconv.ParseFloat(strings.TrimSpace(string(tempData)), 64)
			return math.Round(milli/10) / 100
		}
	}

	return 0
}

func readFileTrim(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}


