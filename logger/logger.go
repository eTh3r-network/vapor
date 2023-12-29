package logger

import "os"
import "log/slog"


func GetLoggers() (*slog.Logger) {
    logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

    return logger
}
