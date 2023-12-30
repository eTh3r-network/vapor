//
// This file is part of the eTh3r project, written, hosted and distributed under MIT License
//  - eTh3r network, 2023-2024
//

package logger

import "log/slog"
import "os"

func GetLogger() (*slog.Logger) {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}
