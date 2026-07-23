package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg time.Time

type themeReloadMsg struct{}

type Metrics struct {
	CPU     CPUMetrics
	GPU     *GPUMetrics
	Memory  MemoryMetrics
	Disks   []DiskMetrics
	Network NetworkMetrics
	Time    time.Time
}

type CPUMetrics struct {
	UsagePercent float64
	Frequency    float64
	Temperature  float64
	Cores        int
}

type GPUMetrics struct {
	UsagePercent float64
	Temperature  float64
	Name         string
	HasMemory    bool
	MemUsed      uint64
	MemTotal     uint64
}

type MemoryMetrics struct {
	UsedBytes      uint64
	TotalBytes     uint64
	AvailableBytes uint64
}

type DiskMetrics struct {
	Mountpoint string
	Fstype     string
	UsedBytes  uint64
	TotalBytes uint64
}

type NetworkMetrics struct {
	RXBytesPerSec float64
	TXBytesPerSec float64
	Interface     string
}

type cpuData struct {
	total uint64
	idle  uint64
}

type netData struct {
	rxBytes uint64
	txBytes uint64
}

type model struct {
	config    Config
	theme     Theme
	metrics   Metrics
	width     int
	height    int
	ready     bool

	themePath    string
	themeModTime time.Time

	prevCPU cpuData
	prevNet map[string]netData
	netTime time.Time
}

func newModel(cfg Config, th Theme) model {
	m := model{
		config:    cfg,
		theme:     th,
		width:     80,
		height:    24,
		themePath: cfg.ThemePath(),
		prevNet:   make(map[string]netData),
	}
	m.prevCPU = readCPUStat()
	m.prevNet = readNetDev()
	m.netTime = time.Now()
	m.collectAbsoluteMetrics()
	return m
}

func (m model) Init() tea.Cmd {
	return m.tickCmd()
}

func (m model) tickCmd() tea.Cmd {
	return tea.Tick(
		time.Duration(m.config.UpdateMs)*time.Millisecond,
		func(t time.Time) tea.Msg { return tickMsg(t) },
	)
}

func (m *model) collectMetrics() {
	m.collectCPUMetrics()
	m.collectNetMetrics()
	m.collectAbsoluteMetrics()
	m.metrics.Time = time.Now()
}

func (m *model) collectAbsoluteMetrics() {
	m.metrics.Memory = readMemory()
	m.metrics.Disks = readDisks()
	m.metrics.GPU = getGPU()
	m.metrics.CPU.Cores = readCPUCores()
	m.metrics.CPU.Frequency = readCPUFreq()
	m.metrics.CPU.Temperature = readCPUTemp()
}

func (m *model) collectCPUMetrics() {
	current := readCPUStat()
	if m.prevCPU.total == 0 {
		m.prevCPU = current
		return
	}
	totalDelta := current.total - m.prevCPU.total
	idleDelta := current.idle - m.prevCPU.idle
	if totalDelta > 0 {
		m.metrics.CPU.UsagePercent = float64(totalDelta-idleDelta) / float64(totalDelta) * 100
	}
	m.prevCPU = current
}

func (m *model) collectNetMetrics() {
	current := readNetDev()
	now := time.Now()
	elapsed := now.Sub(m.netTime).Seconds()
	if elapsed <= 0 {
		m.prevNet = current
		m.netTime = now
		return
	}

	var bestIface string
	var bestSpeed float64

	for name, cur := range current {
		prev, ok := m.prevNet[name]
		if !ok {
			continue
		}
		if cur.rxBytes < prev.rxBytes || cur.txBytes < prev.txBytes {
			continue
		}
		rx := float64(cur.rxBytes-prev.rxBytes) / elapsed
		tx := float64(cur.txBytes-prev.txBytes) / elapsed
		if rx+tx > bestSpeed {
			bestSpeed = rx + tx
			bestIface = name
			m.metrics.Network = NetworkMetrics{
				RXBytesPerSec: rx,
				TXBytesPerSec: tx,
				Interface:     name,
			}
		}
	}

	if bestIface == "" {
		for name := range current {
			if _, ok := m.prevNet[name]; !ok {
				m.prevNet[name] = current[name]
			}
		}
	}

	m.prevNet = current
	m.netTime = now
}
