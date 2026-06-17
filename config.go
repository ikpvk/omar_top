package main

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ColorTheme      string
	UpdateMs        int
	ThemeBackground bool
}

func defaultConfig() Config {
	return Config{
		ColorTheme:      "current",
		UpdateMs:        1000,
		ThemeBackground: true,
	}
}

func (c Config) ThemePath() string {
	home, _ := os.UserHomeDir()
	return home + "/.config/omar_top/themes/" + c.ColorTheme + ".theme"
}

func loadConfig() Config {
	cfg := defaultConfig()
	home, err := os.UserHomeDir()
	if err != nil {
		return cfg
	}
	data, err := os.ReadFile(home + "/.config/omar_top/config.toml")
	if err != nil {
		return cfg
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "[") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		val = strings.Trim(val, "\"'")
		switch key {
		case "color_theme":
			if val != "" {
				cfg.ColorTheme = val
			}
		case "update_ms":
			if ms, err := strconv.Atoi(val); err == nil && ms >= 100 {
				cfg.UpdateMs = ms
			}
		case "theme_background":
			cfg.ThemeBackground = val == "true"
		}
	}
	return cfg
}
