package youtube

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/Microsoft/go-winio"
	"github.com/mahadiksahil60/dev-launcher/internal/paths"
)

const (
	mpvPipe = `\\.\pipe\dev-launcher-mpv`
)

type Player struct {
	cmd *exec.Cmd
}

func NewPlayer() *Player {
	return &Player{}
}

type PlaybackState struct {
	TimePos  float64
	Duration float64
	Paused   bool
}

func (p *Player) Play(url string) error {
	if p.cmd != nil && p.cmd.Process != nil {
		_ = p.cmd.Process.Kill()
		p.cmd = nil
	}

	// NOTE : For local development
	// p.cmd = exec.Command(
	// 	mpvPath,
	// 	"--no-video",
	// 	"--input-ipc-server="+mpvPipe,
	// 	url,
	// )

	p.cmd = exec.Command(
		paths.MPV(),
		"--no-video",
		"--input-ipc-server="+mpvPipe,
		url,
	)

	return p.cmd.Start()
}

func (p *Player) Stop() error {
	if p.cmd == nil || p.cmd.Process == nil {
		return nil
	}

	// Kill the existing MPV process.
	// Ignore the error because the process may have
	// already finished on its own.
	_ = p.cmd.Process.Kill()

	// Clear our reference regardless.
	p.cmd = nil

	return nil
}

func (p *Player) TogglePause() error {
	conn, err := winio.DialPipe(mpvPipe, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	command := map[string]interface{}{
		"command": []interface{}{"cycle", "pause"},
	}

	data, err := json.Marshal(command)
	if err != nil {
		return err
	}

	data = append(data, '\n')

	_, err = conn.Write(data)

	return err
}

func (p *Player) sendCommand(command []interface{}) ([]byte, error) {
	conn, err := winio.DialPipe(mpvPipe, nil)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	request := map[string]interface{}{
		"command": command,
	}

	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	data = append(data, '\n')

	if _, err := conn.Write(data); err != nil {
		return nil, err
	}

	buffer := make([]byte, 4096)

	n, err := conn.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer[:n], nil
}

func (p *Player) GetPlaybackState() (PlaybackState, error) {
	var state PlaybackState

	properties := []string{
		"time-pos",
		"duration",
		"pause",
	}

	for _, property := range properties {
		response, err := p.sendCommand(
			[]interface{}{
				"get_property",
				property,
			},
		)

		if err != nil {
			return state, err
		}

		var result struct {
			Error string          `json:"error"`
			Data  json.RawMessage `json:"data"`
		}

		if err := json.Unmarshal(response, &result); err != nil {
			return state, err
		}

		if result.Error != "success" {
			return state, fmt.Errorf(
				"mpv property %s: %s",
				property,
				result.Error,
			)
		}

		switch property {
		case "time-pos":
			if err := json.Unmarshal(result.Data, &state.TimePos); err != nil {
				return state, err
			}

		case "duration":
			if err := json.Unmarshal(result.Data, &state.Duration); err != nil {
				return state, err
			}

		case "pause":
			if err := json.Unmarshal(result.Data, &state.Paused); err != nil {
				return state, err
			}
		}
	}

	return state, nil
}

func (p *Player) SeekRelative(seconds float64) error {
	_, err := p.sendCommand([]interface{}{
		"seek",
		seconds,
		"relative",
	})

	return err
}
