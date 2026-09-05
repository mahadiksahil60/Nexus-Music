package tui

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m Model) separator() string {
	width := m.width

	if width <= 0 {
		width = 80
	}

	return separatorStyle.Render(strings.Repeat("─", width))
}

func (m Model) renderHeader() string {
	title := titleStyle.Render("N E X U S")
	subtitle := subtitleStyle.Render("MEDIA CONTROL SYSTEM")

	return title + "  " + subtitle
}

func (m Model) renderSystemStatus() string {
	return activeStyle.Render("[ NEXUS CORE ● ENGAGED ]") +
		"  " +
		activeStyle.Render("[ MPV ● CALIBRATED ]") +
		"  " +
		activeStyle.Render("[ YT-ENGINE ● SYNCHRONIZED ]")
}

func (m Model) renderSearch() string {
	prompt := searchPromptStyle.Render("nexus@search:~$")

	input := m.textInput.View()

	return prompt + "  " + input
}

func (m Model) renderResults() string {
	var b strings.Builder

	b.WriteString(sectionStyle.Render("YOUTUBE RESULTS"))
	b.WriteString("\n\n")

	for i, song := range m.songs {
		cursor := " "
		if m.cursor == i {
			cursor = "❯"
		}

		number := fmt.Sprintf("%02d", i+1)

		// Keep the complete row within the terminal.
		maxTitleWidth := m.width - 45

		if maxTitleWidth < 20 {
			maxTitleWidth = 20
		}

		title := truncateText(song.Title, maxTitleWidth)

		if m.cursor == i {
			title = selectedStyle.Render(title)
			cursor = selectedStyle.Render(cursor)
		} else {
			title = normalStyle.Render(title)
		}

		duration := dimStyle.Render(
			fmt.Sprintf("[%s]", song.Duration),
		)

		channel := dimStyle.Render(
			" - " + truncateText(song.Channel, 25),
		)

		line := fmt.Sprintf(
			"%s %s  %s %s %s",
			cursor,
			number,
			title,
			duration,
			channel,
		)

		b.WriteString(line)
		b.WriteString("\n")
	}

	return b.String()
}

func (m Model) renderPlayer() string {
	var b strings.Builder

	b.WriteString(sectionStyle.Render("NOW PLAYING"))
	b.WriteString("\n\n")

	if m.currentlyPlaying == "" {
		b.WriteString(dimStyle.Render("Nothing playing"))
		return b.String()
	}

	frames := []string{
		"▂▅▇▃▆▂▇",
		"▃▇▅█▂▆▃",
		"▇▂▆▃█▅▂",
		"▅█▃▇▂▅█",
	}

	frame := frames[m.playerFrame%len(frames)]

	status := "● PLAYING"

	if m.paused {
		status = "⏸  PAUSED"
		frame = "────────"
	}

	if m.playback_song {
		b.WriteString("  ")
		b.WriteString("On repeat 🔁")
	}

	b.WriteString(
		activeStyle.Render(
			"▶ " + truncateText(m.currentlyPlaying, m.width-20),
		),
	)

	b.WriteString("\n")

	b.WriteString(
		activeStyle.Render(status),
	)

	if !m.paused {
		b.WriteString("  ")
		b.WriteString(activeStyle.Render(frame))
	}

	if m.duration > 0 {
		b.WriteString("\n\n")

		barWidth := m.width - 20

		if barWidth < 20 {
			barWidth = 20
		}

		progress := renderProgress(
			m.currentTime,
			m.duration,
			barWidth,
		)

		b.WriteString(
			dimStyle.Render(formatTime(m.currentTime)),
		)

		b.WriteString(" ")

		b.WriteString(
			activeStyle.Render(progress),
		)

		b.WriteString(" ")

		b.WriteString(
			dimStyle.Render(formatTime(m.duration)),
		)
	}

	b.WriteString("\n")
	b.WriteString(dimStyle.Render("ENGINE: MPV"))

	return b.String()
}

func (m Model) renderCommands() string {
	commands :=
		activeStyle.Render("↑↓") +
			dimStyle.Render(" Navigate") +
			"   " +
			activeStyle.Render("enter") +
			dimStyle.Render(" Play") +
			"   " +
			activeStyle.Render("space") +
			dimStyle.Render(" Pause") +
			"   " +
			activeStyle.Render("/") +
			dimStyle.Render(" Search") +
			"   " +
			activeStyle.Render("q") +
			dimStyle.Render(" Quit") +
			"   " +
			activeStyle.Render("<-<- | ->->") +
			dimStyle.Render(" Fast forward / Fast Backwards")

	credit := activeStyle.Render("◈") + dimStyle.Render(" Forged by Sahil Mahadik 👽")

	return commands + "\n\n" + credit
}

func truncateText(text string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}

	if lipgloss.Width(text) <= maxWidth {
		return text
	}

	runes := []rune(text)

	if len(runes) <= maxWidth {
		return text
	}

	if maxWidth <= 3 {
		return string(runes[:maxWidth])
	}

	return string(runes[:maxWidth-3]) + "..."
}

func scanTick() tea.Cmd {
	return tea.Tick(
		80*time.Millisecond,
		func(time.Time) tea.Msg {
			return scanTickMsg{}
		},
	)
}

type scanTickMsg struct{}

func renderProgress(current, total float64, width int) string {
	if total <= 0 {
		return ""
	}

	if current < 0 {
		current = 0
	}

	if current > total {
		current = total
	}

	progress := current / total

	filled := int(progress * float64(width))

	if filled > width {
		filled = width
	}

	bar := strings.Repeat("━", filled)

	if filled < width {
		bar += "╸"
		bar += strings.Repeat("─", width-filled-1)
	}

	return bar
}

func formatTime(seconds float64) string {
	if seconds < 0 {
		seconds = 0
	}

	minutes := int(seconds) / 60
	secondsPart := int(seconds) % 60

	return fmt.Sprintf(
		"%02d:%02d",
		minutes,
		secondsPart,
	)
}
