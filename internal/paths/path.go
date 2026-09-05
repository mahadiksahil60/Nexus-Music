package paths

import (
	"fmt"
	"os"
	"path/filepath"
)

func BaseDir() string {

	// First, try the directory containing the executable.
	if exe, err := os.Executable(); err == nil {

		exeDir := filepath.Dir(exe)

		fmt.Println(exeDir, "current file path.")

		mpvPath := filepath.Join(
			exeDir,
			"bin",
			"mpv",
			"mpv.com",
		)

		ytDlpPath := filepath.Join(
			exeDir,
			"bin",
			"yt-dlp",
			"yt-dlp.exe",
		)

		// Portable/release layout.
		if fileExists(mpvPath) && fileExists(ytDlpPath) {
			return exeDir
		}
	}

	// Development mode:
	// go run creates a temporary executable,
	// so use the project directory.
	if cwd, err := os.Getwd(); err == nil {
		return cwd
	}

	return "."
}

func MPV() string {
	return filepath.Join(
		BaseDir(),
		"bin",
		"mpv",
		"mpv.com",
	)
}

func YTDLP() string {
	return filepath.Join(
		BaseDir(),
		"bin",
		"yt-dlp",
		"yt-dlp.exe",
	)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
