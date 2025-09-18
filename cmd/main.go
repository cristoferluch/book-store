package main

import (
	"book-store/configs"
	"book-store/pkg/database"
	"book-store/pkg/logger"
	"book-store/pkg/server"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	if err := start(); err != nil {
		log.Fatalf("Error starting server: %v", err.Error())
	}

}

func start() error {

	err := configs.NewConfigLoader()
	if err != nil {
		return err
	}

	lg, err := logger.NewZapLogger()
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
		Addr:         fmt.Sprintf(":%d", configs.Cfg.App.Port),
		Handler:      r,
		ReadTimeout:  time.Second * configs.Cfg.App.ReadTimeout,
		WriteTimeout: time.Second * configs.Cfg.App.WriteTimeout,
		IdleTimeout:  time.Second * configs.Cfg.App.IdleTimeout,
	}

	go func() {
		zap.L().Info("server started",
			zap.String("addr", fmt.Sprintf(":%d", configs.Cfg.App.Port)),
			zap.String("version", configs.Cfg.App.Version),
			zap.String("environment", configs.Cfg.App.Environment),
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
