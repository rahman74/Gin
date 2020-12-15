package utils

import (
	"log"
	"test/model"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// InitializeLogSetting ...
func InitializeLogSetting(logConfig *model.LogConfig) {
	if logConfig.Filename != "" {
		log.SetOutput(
			&lumberjack.Logger{
				Filename:   logConfig.Filename,
				MaxSize:    logConfig.MaxSize,
				MaxBackups: logConfig.MaxBackups,
				MaxAge:     logConfig.MaxAge,
				Compress:   logConfig.Compress,
			},
		)
	}

}
