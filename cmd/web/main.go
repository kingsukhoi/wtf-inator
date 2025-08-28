package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/kingsukhoi/wtf-inator/pkg/conf"
	"github.com/kingsukhoi/wtf-inator/pkg/db"
	"github.com/kingsukhoi/wtf-inator/pkg/proxy"
	"github.com/kingsukhoi/wtf-inator/pkg/routes"
	slogecho "github.com/samber/slog-echo"
)

func main() {

	config := conf.MustGetConfig("config.yaml")

	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelDebug)

	var logger *slog.Logger
	if config.JsonLogs {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))
		slog.SetDefault(logger)
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))
		slog.SetDefault(logger)
	}

	//start up db pool
	_ = db.MustGetDatabase()

	e, err := routes.NewRouter()
	if err != nil {
		panic(err)
	}

	slogConfig := slogecho.Config{
		WithSpanID:    true,
		WithTraceID:   true,
		WithRequestID: true,
	}

	e.Use(slogecho.NewWithConfig(logger, slogConfig))

	osInterruptContext, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		echoErr := e.Start(config.Port)
		if echoErr != nil {
			if errors.Is(echoErr, http.ErrServerClosed) {
				slog.Info("shutting down the server")
			} else {
				slog.Error("error with server", "exception", echoErr)
			}
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	<-osInterruptContext.Done()
	cleanupContext, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	err = e.Shutdown(cleanupContext)
	if err != nil {
		slog.Error("error with server shutdown", "exception", err)
	}

	proxy.Wait()
}
