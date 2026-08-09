package tui

import (
	"fmt"
	"os"

	"github.com/mahadiksahil60/dev-launcher/internal/paths"
)

type DependencyStatus struct {
	MPV   bool
	YTDLP bool
}

func CheckDependencies() DependencyStatus {
	return DependencyStatus{
		MPV:   fileExists(paths.MPV()),
		YTDLP: fileExists(paths.YTDLP()),
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func DependencyError(status DependencyStatus) error {
	var missing string

	if !status.MPV {
		missing += fmt.Sprintf(
			"MPV engine missing:\n  %s\n\n",
			paths.MPV(),
		)
	}

	if !status.YTDLP {
		missing += fmt.Sprintf(
			"yt-dlp missing:\n  %s\n\n",
			paths.YTDLP(),
		)
	}

	if missing == "" {
		return nil
	}

	return fmt.Errorf(
		"NEXUS runtime dependency check failed:\n\n%s",
		missing,
	)
}
