package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if !m.ready {
			m.ready = true
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}

	case tickMsg:
		if th, modTime, ok := checkThemeReload(m.themePath, m.themeModTime); ok {
			m.theme = th
			m.themeModTime = modTime
		}
		m.collectMetrics()
		return m, m.tickCmd()

	case themeReloadMsg:
		th, modTime, ok := checkThemeReload(m.themePath, m.themeModTime)
		if ok {
			m.theme = th
			m.themeModTime = modTime
		} else {
			th := loadTheme(m.themePath)
			m.theme = th
			m.themeModTime = modTime
		}
	}

	return m, nil
}
