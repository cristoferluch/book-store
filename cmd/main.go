package main

import (
	"book-store/pkg/config"
	"book-store/pkg/server"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"book-store/pkg/database"
	"syscall"
	"time"

	"github.com/spf13/viper"
)

func main() {

	if err := start(); err != nil {
		slog.Error("Error starting server", slog.String("error", err.Error()))
		panic(err)
	}
}

func start() error {

	cfg, err := newConfigLoader()
	if err != nil {
		return err
	}

	if err := setupLogger(); err != nil {
		return err
	}

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	r := server.NewServer(db)

	s := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.App.Port),
		Handler:      r,
		ReadTimeout:  time.Second * cfg.App.ReadTimeout,
		WriteTimeout: time.Second * cfg.App.WriteTimeout,
		IdleTimeout:  time.Second * cfg.App.IdleTimeout,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("server started",
			slog.String("addr", fmt.Sprintf(":%d", cfg.App.Port)),
			slog.String("version", cfg.App.Version),
		)
		if err := s.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
		}
	}()

	<-quit
	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error",
				slog.String("error", err.Error()),
			)
		}
	}

	slog.Info("Server successfully shut down")

	return nil
}

func setupLogger() error {

	err := os.Mkdir("./logs", 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	fileName := fmt.Sprintf("./logs/%s.log", time.Now().Format("2006-01-02"))

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	multiWriter := io.MultiWriter(os.Stdout, file)

	logger := slog.New(slog.NewJSONHandler(multiWriter, nil))

	slog.SetDefault(logger)

	return nil
}

func newConfigLoader() (*config.Config, error) {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var cfg config.Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	return &cfg, nil
}
