package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	leftMargin  = 2
	labelWidth  = 6
	suffixWidth = 5
)

func barPrefixWidth() int {
	return leftMargin + labelWidth + 1
}

func barWidth(totalWidth int) int {
	w := totalWidth - barPrefixWidth() - 2 - suffixWidth
	if w < 10 {
		w = 10
	}
	return w
}

func (m model) View() string {
	if !m.ready {
		return ""
	}
	var b strings.Builder

	b.WriteString(m.renderHeader())
	b.WriteString("\n")
	b.WriteString(m.renderDivider())
	b.WriteString("\n")

	b.WriteString(m.renderCPUSection())
	b.WriteString("\n")

	if m.metrics.GPU != nil {
		b.WriteString(m.renderGPUSection())
		b.WriteString("\n")
	}

	b.WriteString(m.renderMemSection())
	b.WriteString("\n")

	b.WriteString(m.renderDiskSection())
	b.WriteString("\n")

	b.WriteString(m.renderNetSection())

	content := b.String()
	lines := strings.Split(content, "\n")

	var result strings.Builder
	for i, line := range lines {
		if i == len(lines)-1 && line == "" {
			continue
		}
		style := lipgloss.NewStyle().
			Foreground(m.theme.Fg).
			Width(m.width)
		if m.config.ThemeBackground {
			style = style.Background(m.theme.Bg)
		}
		rendered := style.Render(line)
		result.WriteString(rendered)
		if i < len(lines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

func (m model) renderHeader() string {
	title := lipgloss.NewStyle().
		Foreground(m.theme.Title).
		Bold(true).
		Render("OMAR-TOP")

	refresh := formatRefresh(m.config.UpdateMs)

	header := fmt.Sprintf("%s%s%s",
		strings.Repeat(" ", leftMargin),
		title,
		lipgloss.NewStyle().Foreground(m.theme.Label).Render("  "+refresh),
	)

	return header
}

func formatRefresh(ms int) string {
	if ms < 1000 {
		return fmt.Sprintf("%dms refresh", ms)
	}
	return fmt.Sprintf("%ds refresh", ms/1000)
}

func (m model) renderDivider() string {
	div := strings.Repeat("─", m.width-leftMargin*2)
	return strings.Repeat(" ", leftMargin) +
		lipgloss.NewStyle().Foreground(m.theme.Divider).Render(div)
}

func (m model) renderLabel(label string) string {
	padded := fmt.Sprintf("%-*s", labelWidth, label)
	return strings.Repeat(" ", leftMargin) +
		lipgloss.NewStyle().Foreground(m.theme.Label).Render(padded) + " "
}

func (m model) renderPct(pct float64) string {
	return lipgloss.NewStyle().
		Foreground(m.theme.Value).
		Render(fmt.Sprintf(" %3.0f%%", pct))
}

func (m model) renderBar(pct float64) string {
	bw := barWidth(m.width)
	filled := int(math.Round(pct * float64(bw) / 100.0))
	if filled < 0 {
		filled = 0
	}
	if filled > bw {
		filled = bw
	}
	fillStyle := lipgloss.NewStyle().Foreground(m.theme.BarFill)
	emptyStyle := lipgloss.NewStyle().Foreground(m.theme.BarBg)
	return fillStyle.Render(strings.Repeat("█", filled)) +
		emptyStyle.Render(strings.Repeat("░", bw-filled))
}

func (m model) renderInfo(text string) string {
	indent := barPrefixWidth()
	return strings.Repeat(" ", indent) +
		lipgloss.NewStyle().Foreground(m.theme.Info).Render(text)
}

func (m model) renderSectionLine(label string, pct float64) string {
	return m.renderLabel(label) + m.renderBar(pct) + m.renderPct(pct)
}

func (m model) renderCPUSection() string {
	cpu := m.metrics.CPU
	label := "CPU"
	line1 := m.renderSectionLine(label, cpu.UsagePercent)

	var infoParts []string
	if cpu.Frequency > 0 {
		infoParts = append(infoParts, formatFreq(cpu.Frequency))
	}
	if cpu.Temperature > 0 {
		infoParts = append(infoParts, formatTemp(cpu.Temperature))
	}
	if len(infoParts) == 0 {
		return line1
	}
	line2 := m.renderInfo(strings.Join(infoParts, " · "))
	return line1 + "\n" + line2
}

func (m model) renderGPUSection() string {
	gpu := m.metrics.GPU
	if gpu == nil {
		return ""
	}
	name := gpu.Name
	if len(name) > 12 {
		name = name[:12]
	}
	label := fmt.Sprintf("%-6s", name)

	line1 := m.renderLabel(label) + m.renderBar(gpu.UsagePercent) + m.renderPct(gpu.UsagePercent)

	var infoParts []string
	if gpu.Temperature > 0 {
		infoParts = append(infoParts, formatTemp(gpu.Temperature))
	}
	if gpu.HasMemory && gpu.MemTotal > 0 {
		pct := float64(gpu.MemUsed) / float64(gpu.MemTotal) * 100
		infoParts = append(infoParts, fmt.Sprintf("%s / %s (%.0f%%)",
			formatBytes(gpu.MemUsed), formatBytes(gpu.MemTotal), pct))
	}
	if len(infoParts) == 0 {
		return line1
	}
	line2 := m.renderInfo(strings.Join(infoParts, " · "))
	return line1 + "\n" + line2
}

func (m model) renderMemSection() string {
	mem := m.metrics.Memory
	label := "MEM"
	if mem.TotalBytes == 0 {
		return m.renderSectionLine(label, 0)
	}
	pct := float64(mem.UsedBytes) / float64(mem.TotalBytes) * 100
	line1 := m.renderSectionLine(label, pct)
	availPct := float64(mem.AvailableBytes) / float64(mem.TotalBytes) * 100
	line2 := m.renderInfo(fmt.Sprintf("%s / %s  (%.0f%% avail)",
		formatBytes(mem.UsedBytes), formatBytes(mem.TotalBytes), availPct))
	return line1 + "\n" + line2
}

func (m model) renderDiskSection() string {
	disks := m.metrics.Disks
	if len(disks) == 0 {
		return m.renderSectionLine("DISK", 0)
	}

	var lines []string
	for i, disk := range disks {
		mp := "/"
		if disk.Mountpoint != "/" {
			mp = disk.Mountpoint
		}
		if len(mp) > 5 {
			mp = mp[:5]
		}
		label := fmt.Sprintf("%-6s", "DSK "+mp)
		if disk.TotalBytes == 0 {
			lines = append(lines, m.renderSectionLine(label, 0))
		} else {
			pct := float64(disk.UsedBytes) / float64(disk.TotalBytes) * 100
			line1 := m.renderSectionLine(label, pct)
			line2 := m.renderInfo(fmt.Sprintf("%s / %s",
				formatBytes(disk.UsedBytes), formatBytes(disk.TotalBytes)))
			lines = append(lines, line1+"\n"+line2)
		}
		if i == 0 {
			break
		}
	}
	return strings.Join(lines, "\n")
}

func (m model) renderNetSection() string {
	label := "NET"
	net := m.metrics.Network
	if net.Interface == "" {
		return m.renderSectionLine(label, 0)
	}
	line1 := m.renderLabel(label) + m.renderNetText(net)

	total := net.RXBytesPerSec + net.TXBytesPerSec
	maxSpeed := total * 1.5
	if maxSpeed < 1024*1024 {
		maxSpeed = 1024 * 1024
	}
	pct := total / maxSpeed * 100
	if pct > 100 {
		pct = 100
	}

	line2 := m.renderInfo(fmt.Sprintf("%s on %s",
		formatSpeed(total), net.Interface))
	return line1 + "\n" + line2
}

func (m model) renderNetText(net NetworkMetrics) string {
	bw := barWidth(m.width)
	down := formatSpeed(net.RXBytesPerSec)
	up := formatSpeed(net.TXBytesPerSec)
	text := fmt.Sprintf("↓ %s  ↑ %s", down, up)
	if len(text) > bw {
		text = fmt.Sprintf("↓%s ↑%s", down[:min(len(down), 8)], up[:min(len(up), 8)])
	}
	return lipgloss.NewStyle().Foreground(m.theme.Value).Render(
		fmt.Sprintf("%-*s", bw, text))
}
