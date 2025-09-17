package main

import (
	"book-store/pkg/config"
	"book-store/pkg/database"
	"book-store/pkg/server"
	"context"
	"errors"
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/spf13/viper"
)

func main() {

	if err := start(); err != nil {
		log.Fatalf("Error starting server: %v", err.Error())
	}

}

func start() error {

	err := newConfigLoader()
	if err != nil {
		return err
	}

	lg, err := newZapLogger()
	if err != nil {
		return err
	}
	defer lg.Sync()

	undo := zap.ReplaceGlobals(lg)
	defer undo()

	db, err := database.NewPostgresDB()
	if err != nil {
		return err
	}
	defer db.Close()

	r := server.NewServer(db)

	s := http.Server{
		Addr:         fmt.Sprintf(":%d", config.Cfg.App.Port),
		Handler:      r,
		ReadTimeout:  time.Second * config.Cfg.App.ReadTimeout,
		WriteTimeout: time.Second * config.Cfg.App.WriteTimeout,
		IdleTimeout:  time.Second * config.Cfg.App.IdleTimeout,
	}

	go func() {
		zap.L().Info("server started",
			zap.String("addr", fmt.Sprintf(":%d", config.Cfg.App.Port)),
			zap.String("version", config.Cfg.App.Version),
			zap.String("environment", config.Cfg.App.Environment),
		)
		if err := s.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	zap.L().Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			zap.L().Error("server error",
				zap.String("error", err.Error()),
			)
		}
	}

	zap.L().Info("server successfully shut down")

	return nil
}

func newZapLogger() (logger *zap.Logger, err error) {

	if config.Cfg.App.Environment == "production" {
		if err := os.MkdirAll(config.Cfg.App.LogPath, 0755); err != nil {
			return nil, err
		}
	}

	var encoderConfig zapcore.EncoderConfig
	var level zapcore.Level

	if config.Cfg.App.Environment == "production" {
		encoderConfig = zap.NewProductionEncoderConfig()
		level = zapcore.InfoLevel
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		level = zapcore.DebugLevel
	}

	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	fileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   path.Join(config.Cfg.App.LogPath, "app.log"),
		MaxSize:    10,
		MaxBackups: 7,
		MaxAge:     30,
		Compress:   false,
	})

	var core zapcore.Core

	if config.Cfg.App.Environment == "production" {
		core = zapcore.NewCore(encoder, fileWriteSyncer, level)
	} else {
		core = zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(os.Stdout), level)
	}

	return zap.New(core, zap.AddCaller()), nil
}

func newConfigLoader() error {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("error reading config file, %w", err)
	}

	err = viper.Unmarshal(&config.Cfg)
	if err != nil {
		return fmt.Errorf("unable to decode into struct, %w", err)
	}

	config.Cfg.StartTime = time.Now()

	return nil
}
