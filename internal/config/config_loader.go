package config

import "github.com/atomreforge/confx"

// MustLoadConfig 加载配置，出错时 panic。
func MustLoadConfig(paths ...string) *Config {
	opts := []confx.Option{
		confx.WithFileName("config"),
		confx.WithFileType("yaml"),
	}
	if len(paths) > 0 {
		opts = append(opts, confx.WithSearchPaths(paths...))
	}

	// 找不到配置文件时用默认值，不计为错误（便于测试环境下直接运行）。
	cfg := confx.MustLoad[Config]("dnappcli", opts...)
	return cfg
}
