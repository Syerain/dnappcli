package utils

import (
	"log/slog"
	"os"
	"time"

	"github.com/pwntr/tinter"
)

// InitModuleLogger 初始化带颜色、模块名与 debug 级别的 slog 日志器，
// 用法与 server 端 internal/utils/logger.go 保持一致。
func InitModuleLogger(isDebugMode bool, moduleName string) {
	var handler slog.Handler

	level := slog.LevelInfo
	if isDebugMode {
		level = slog.LevelDebug
	}

	handler = tinter.NewHandler(os.Stdout, &tinter.Options{
		Level:      level,
		TimeFormat: time.RFC3339,
	})

	logger := slog.New(handler).With(
		slog.String("module", moduleName),
	)

	slog.SetDefault(logger)
}

func GetLogger() *slog.Logger {
	return slog.Default()
}
