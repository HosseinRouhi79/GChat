package logging

import "chat/config"

type Logger interface {
	Init()

	Debug(cat Category, subCat SubCategory, msg string, extra map[ExtraKey]interface{})
	Debugf(tempplate string, args ...interface{})

	Info(cat Category, subCat SubCategory, msg string, extra map[ExtraKey]interface{})
	Infof(tempplate string, args ...interface{})

	Fatal(cat Category, subCat SubCategory, msg string, extra map[ExtraKey]interface{})
	Fatalf(tempplate string, args ...interface{})

	Warn(cat Category, subCat SubCategory, msg string, extra map[ExtraKey]interface{})
	Warnf(tempplate string, args ...interface{})
}

func NewLogger(cfg *config.Config) Logger {
	if cfg.Logger.Logger == "zap" {
		return newZapLogger(cfg)
	} else if cfg.Logger.Logger == "zerolog" {
		return newZeroLogger(cfg)
	}
	return nil
}
