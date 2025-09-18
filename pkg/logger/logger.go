package logger

import (
	"book-store/configs"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"time"
)

func NewZapLogger() (logger *zap.Logger, err error) {

	switch strings.ToLower(configs.Cfg.App.Environment) {

	case "production", "prod":

		file, err := createLogFile()
		if err != nil {
			return nil, err
		}

		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		encoder := zapcore.NewJSONEncoder(encoderConfig)
		multiSync := zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(file))

		return zap.New(
			zapcore.NewCore(encoder, multiSync, zapcore.InfoLevel),
			zap.AddCaller(),
		), nil

	case "development", "dev":

		encoderConfig := zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder := zapcore.NewConsoleEncoder(encoderConfig)

		return zap.New(
			zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
			zap.AddCaller(),
		), nil
	}

	return nil, errors.New("unknown environment: " + configs.Cfg.App.Environment)
}

func createLogFile() (*os.File, error) {

	fileName := fmt.Sprintf("%s/%s.log", configs.Cfg.App.LogPath, time.Now().Format("20060102"))

	if err := os.MkdirAll(configs.Cfg.App.LogPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("error creating log directory: %w", err)
	}

	file, err := os.Create(fileName)
	if err != nil {
		return nil, fmt.Errorf("error creating log file: %w", err)
	}

	return file, nil
}
