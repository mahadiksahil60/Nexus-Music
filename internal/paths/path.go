package paths

import (
	"os"
	"path/filepath"
)

func ExecutableDir() string {
	exe, err := os.Executable()
	if err != nil {
		return "."
	}

	return filepath.Dir(exe)
}

func MPV() string {
	return filepath.Join(
		ExecutableDir(),
		"bin",
		"mpv",
		"mpv.exe",
	)
}

func YTDLP() string {
	return filepath.Join(
		ExecutableDir(),
		"bin",
		"yt-dlp",
		"yt-dlp.exe",
	)
}
