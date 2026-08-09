package tui

import "charm.land/lipgloss/v2"

var (
	primaryColor = lipgloss.Color("#00FF66")
	cyanColor    = lipgloss.Color("#00FFFF")
	textColor    = lipgloss.Color("#D0D0D0")
	mutedColor   = lipgloss.Color("#777777")
	errorColor   = lipgloss.Color("#FF4444")
	warningColor = lipgloss.Color("#FFD700")

	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	sectionStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	selectedStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(textColor)

	dimStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	activeStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	cyanStyle = lipgloss.NewStyle().
			Foreground(cyanColor)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	warningStyle = lipgloss.NewStyle().
			Foreground(warningColor)

	separatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#333333"))

	searchStyle = lipgloss.NewStyle().
			Foreground(textColor)

	searchPromptStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true)
)
