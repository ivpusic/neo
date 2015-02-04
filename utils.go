package neo

import (
	"errors"
	"github.com/ivpusic/golog"
	"os"
	"strings"
)

// Checking if provided path is directory
func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

// Merging multiple ``appliable`` slices into one.
func merge(appliabels ...[]appliable) []appliable {
	all := []appliable{}

	for _, appl := range appliabels {
		all = append(all, appl...)
	}

	return all
}

// For given string representation of log level, return ``golog.Level`` instance.
func parseLogLevel(level string) (golog.Level, error) {
	level = strings.ToLower(level)

	switch level {
	case "debug":
		return golog.DEBUG, nil
	case "info":
		return golog.INFO, nil
	case "warn":
		return golog.WARN, nil
	case "error":
		return golog.ERROR, nil
	case "panic":
		return golog.PANIC, nil
	}

	return golog.Level{}, errors.New("Log level " + level + " not found!")
}
