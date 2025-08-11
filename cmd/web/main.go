package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kingsukhoi/wtf-inator/pkg/conf"
	"github.com/kingsukhoi/wtf-inator/pkg/proxy"
	"github.com/kingsukhoi/wtf-inator/pkg/routes"
)

func main() {

	config := conf.MustGetConfig("config.yaml")

	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelDebug)

	if config.JsonLogs {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))
		slog.SetDefault(logger)
	} else {
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))
		slog.SetDefault(logger)
	}

	e, err := routes.NewRouter()
	if err != nil {
		panic(err)
	}

	e.Logger.Fatal(e.Start(":1323"))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	<-ctx.Done()

	proxy.Wait()
}
