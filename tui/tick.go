package tui

import (
	"time"

	tea "charm.land/bubbletea/v2"
)

type uiTickMsg struct{}

func uiTick() tea.Cmd {
	return tea.Tick(
		500*time.Millisecond,
		func(time.Time) tea.Msg {
			return uiTickMsg{}
		},
	)
}
