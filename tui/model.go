package tui

import (
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/mahadiksahil60/dev-launcher/youtube"
)

// state memory of our UI
type Model struct {
	songs            []youtube.Song
	cursor           int
	player           *youtube.Player
	currentlyPlaying string

	query   string
	loading bool
	err     error

	searching bool
	textInput textinput.Model

	// terminal dimensions
	width  int
	height int

	// progress bars
	scanProgress int
	scanFrame    int

	// animation
	playerFrame int

	// playback progress bar.
	currentTime float64
	duration    float64
	paused      bool

	// fast forwand
	lastArrow     string
	lastArrowTime time.Time

	// repeat playback
	playback_song bool
}

type searchResultMsg struct {
	songs []youtube.Song
	err   error
}

func NewModel() Model {
	ti := textinput.New()
	ti.Prompt = "Search: "
	ti.Focus()

	player := youtube.NewPlayer()

	if err := player.Start(); err != nil {
		fmt.Println("Failed to start MPV:", err)
	}

	return Model{
		textInput: ti,
		player:    player,
		searching: true,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		uiTick(),
	)
}

func searchCommand(query string) tea.Cmd {
	return func() tea.Msg {
		songs, err := youtube.Search(query)

		return searchResultMsg{
			songs: songs,
			err:   err,
		}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case uiTickMsg:
		if m.currentlyPlaying != "" {
			m.playerFrame++

			state, err := m.player.GetPlaybackState()

			if err == nil {
				m.currentTime = state.TimePos
				m.duration = state.Duration
				m.paused = state.Paused
			}
		}

		return m, uiTick()

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		return m, nil

	case scanTickMsg:
		if !m.loading {
			return m, nil
		}

		m.scanFrame++

		if m.scanProgress < 90 {
			m.scanProgress += 2
		}

		return m, scanTick()

	case tea.KeyPressMsg:

		// Handle keyboard input differently while searching.
		if m.searching {
			switch msg.String() {

			case "enter":
				m.searching = false
				m.textInput.Blur()

				m.query = m.textInput.Value()
				m.loading = true
				m.err = nil

				m.scanProgress = 0
				m.scanFrame = 0

				return m, tea.Batch(
					searchCommand(m.query),
					scanTick(),
				)

			case "esc":
				m.searching = false
				m.textInput.Blur()

				return m, nil
			}

			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)

			return m, cmd
		}

		// Normal browsing mode.
		switch msg.String() {

		case "q", "ctrl+c":
			m.player.Stop()
			return m, tea.Quit

		case "/":
			m.searching = true
			m.textInput.SetValue("")
			m.textInput.Focus()

			return m, textinput.Blink

		case "up":
			if len(m.songs) == 0 {
				return m, nil
			}

			m.cursor--

			if m.cursor < 0 {
				m.cursor = len(m.songs) - 1
			}

		case "down":
			if len(m.songs) == 0 {
				return m, nil
			}

			m.cursor++

			if m.cursor >= len(m.songs) {
				m.cursor = 0
			}
		case "enter":
			if len(m.songs) == 0 {
				return m, nil
			}

			song := m.songs[m.cursor]

			if err := m.player.Play(song.URL); err != nil {
				m.err = err
				return m, nil
			}

			m.currentlyPlaying = song.Title
			m.currentTime = 0
			m.duration = 0
			m.paused = false

		case "p":
			m.playback_song = !m.playback_song
			for m.playback_song {

			}

		case "left":
			if m.lastArrow == "left" &&
				time.Since(m.lastArrowTime) < 400*time.Millisecond {

				if err := m.player.SeekRelative(-10); err != nil {
					m.err = err
				}

				m.lastArrow = ""
				return m, nil
			}

			m.lastArrow = "left"
			m.lastArrowTime = time.Now()

		case "right":
			if m.lastArrow == "right" &&
				time.Since(m.lastArrowTime) < 400*time.Millisecond {

				if err := m.player.SeekRelative(10); err != nil {
					m.err = err
				}

				m.lastArrow = ""
				return m, nil
			}

			m.lastArrow = "right"
			m.lastArrowTime = time.Now()
		case "space":
			if err := m.player.TogglePause(); err != nil {
				m.err = err
			}
		}

	case searchResultMsg:
		m.loading = false
		m.songs = msg.songs
		m.err = msg.err
		m.cursor = 0
		m.scanProgress = 100

		return m, nil
	}

	return m, nil
}

func (m Model) renderScanner() string {
	barWidth := m.width - 20

	if barWidth < 20 {
		barWidth = 20
	}

	filled := barWidth * m.scanProgress / 100

	if filled > barWidth {
		filled = barWidth
	}

	bar := strings.Repeat("█", filled)
	bar += strings.Repeat("░", barWidth-filled)

	frames := []string{
		"resolving query",
		"contacting youtube",
		"fetching results",
		"parsing metadata",
	}

	frame := frames[m.scanFrame%len(frames)]

	return warningStyle.Render("[ SCANNING YOUTUBE ]") +
		"\n\n" +
		activeStyle.Render(bar) +
		fmt.Sprintf("  %d%%", m.scanProgress) +
		"\n" +
		dimStyle.Render("> "+frame)
}

func (m Model) View() tea.View {
	var b strings.Builder

	// Header
	b.WriteString(m.renderHeader())
	b.WriteString("\n")
	b.WriteString(m.renderSystemStatus())
	b.WriteString("\n")
	b.WriteString(m.separator())
	b.WriteString("\n\n")

	// Search mode
	if m.searching {
		b.WriteString(m.renderSearch())
		b.WriteString("\n\n")
		b.WriteString(
			dimStyle.Render("Press Enter to search   Esc to cancel"),
		)

		return tea.NewView(b.String())
	}

	// Loading
	if m.loading {
		if m.loading {
			b.WriteString(m.renderScanner())
			b.WriteString("\n\n")
			b.WriteString(m.renderCommands())

			return tea.NewView(b.String())
		}
	}

	// Error
	if m.err != nil {
		b.WriteString(
			errorStyle.Render(
				fmt.Sprintf("[ ERROR ] %v", m.err),
			),
		)

		b.WriteString("\n\n")
		b.WriteString(m.renderCommands())

		return tea.NewView(b.String())
	}

	// Search prompt
	b.WriteString(
		dimStyle.Render("nexus@search:~$ "),
	)

	b.WriteString(
		normalStyle.Render(m.query),
	)

	b.WriteString("\n\n")

	// Results
	b.WriteString(m.renderResults())

	b.WriteString("\n")
	b.WriteString(m.separator())
	b.WriteString("\n\n")

	// Player
	b.WriteString(m.renderPlayer())

	b.WriteString("\n\n")
	b.WriteString(m.separator())
	b.WriteString("\n\n")

	// Commands
	b.WriteString(renderSystemStats())
	b.WriteString("\n")
	b.WriteString(m.renderCommands())

	return tea.NewView(b.String())
}
