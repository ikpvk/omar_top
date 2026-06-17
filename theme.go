package main

import (
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	Bg      lipgloss.Color
	Fg      lipgloss.Color
	Title   lipgloss.Color
	Accent  lipgloss.Color
	BarBg   lipgloss.Color
	BarFill lipgloss.Color
	Label   lipgloss.Color
	Value   lipgloss.Color
	Divider lipgloss.Color
	Info    lipgloss.Color
}

func defaultTheme() Theme {
	return Theme{
		Bg:      lipgloss.Color("#1e1e2e"),
		Fg:      lipgloss.Color("#cdd6f4"),
		Title:   lipgloss.Color("#cdd6f4"),
		Accent:  lipgloss.Color("#89b4fa"),
		BarBg:   lipgloss.Color("#585b70"),
		BarFill: lipgloss.Color("#89b4fa"),
		Label:   lipgloss.Color("#585b70"),
		Value:   lipgloss.Color("#cdd6f4"),
		Divider: lipgloss.Color("#585b70"),
		Info:    lipgloss.Color("#a6adc8"),
	}
}

func loadTheme(path string) Theme {
	th := defaultTheme()
	data, err := os.ReadFile(path)
	if err != nil {
		return th
	}
	colors := map[string]*lipgloss.Color{
		"bg":       &th.Bg,
		"fg":       &th.Fg,
		"title":    &th.Title,
		"accent":   &th.Accent,
		"bar_bg":   &th.BarBg,
		"bar_fill": &th.BarFill,
		"label":    &th.Label,
		"value":    &th.Value,
		"divider":  &th.Divider,
		"info":     &th.Info,
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if !strings.HasPrefix(line, "theme[") {
			continue
		}
		bracketEnd := strings.Index(line, "]")
		if bracketEnd < 0 {
			continue
		}
		key := line[6:bracketEnd]
		rest := line[bracketEnd+1:]
		eqIdx := strings.Index(rest, "=")
		if eqIdx < 0 {
			continue
		}
		val := strings.TrimSpace(rest[eqIdx+1:])
		val = strings.Trim(val, "\"")
		if ptr, ok := colors[key]; ok && strings.HasPrefix(val, "#") {
			*ptr = lipgloss.Color(val)
		}
	}
	return th
}

func checkThemeReload(path string, modTime time.Time) (Theme, time.Time, bool) {
	info, err := os.Stat(path)
	if err != nil {
		return Theme{}, modTime, false
	}
	newMod := info.ModTime()
	if newMod.Equal(modTime) {
		return Theme{}, modTime, false
	}
	return loadTheme(path), newMod, true
}
