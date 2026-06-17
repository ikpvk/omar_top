package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cfg := loadConfig()
	th := loadTheme(cfg.ThemePath())

	m := newModel(cfg, th)
	m.themeModTime = getModTime(cfg.ThemePath())

	p := tea.NewProgram(m, tea.WithAltScreen())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGUSR2)
	go func() {
		for range sigCh {
			p.Send(themeReloadMsg{})
		}
	}()

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func getModTime(path string) time.Time {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}
	}
	return info.ModTime()
}
