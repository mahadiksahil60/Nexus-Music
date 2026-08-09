package youtube

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/mahadiksahil60/dev-launcher/internal/paths"
)

// creating custom type for song
type Song struct {
	ID       string
	Title    string
	Channel  string
	Duration string
	URL      string
}

type ytDlpResult struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Channel    string `json:"channel"`
	Duration   string `json:"duration_string"`
	WebpageURL string `json:"webpage_url"`
}

// function to search songs
func Search(query string) ([]Song, error) {

	// NOTE : For local development.
	// cmd := exec.Command(
	// 	"yt-dlp",
	// 	"ytsearch10:"+query,
	// 	"--dump-json",
	// 	"--flat-playlist",
	// 	"--no-warnings",
	// )

	cmd := exec.Command(
		paths.YTDLP(),
		"ytsearch10:"+query,
		"--dump-json",
		"--flat-playlist",
		"--no-warnings",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		fmt.Print(err)
		return nil, err
	}

	// Storing search results.
	var songs []Song

	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		var result ytDlpResult

		if err := json.Unmarshal(scanner.Bytes(), &result); err != nil {
			continue
		}

		songs = append(songs, Song{
			ID:       result.ID,
			Title:    result.Title,
			Channel:  result.Channel,
			Duration: result.Duration,
			URL:      result.WebpageURL,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("yt-dlp failed: %w", err)
	}

	return songs, nil
}
