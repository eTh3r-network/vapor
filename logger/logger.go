package logger

import "log/slog"
import "os"

func GetLogger() (*slog.Logger) {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}
