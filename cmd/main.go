package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Novochenko/sso/internal/app"
	"github.com/Novochenko/sso/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	fmt.Println(cfg)
	dbURL := DBUrlSetup(cfg)
	application := app.New(log, cfg.GRPC.Port, dbURL, cfg.TokenTTL)

	go application.GRPCServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func DBUrlSetup(cfg *config.Config) (DBUrl string) {
	// DBUrl = fmt.Sprintf("%s:%s@tcp(%s)/%s?multiStatements=true&parseTime=true", //multiStatements=true&
	// 	cfg.StoragePath.User,
	// 	cfg.StoragePath.Password,
	// 	cfg.StoragePath.Host,
	// 	cfg.StoragePath.DBName)
	DBUrl = fmt.Sprintf("%s:%s@tcp(%s)/%s?multiStatements=true&parseTime=true", //multiStatements=true&
		cfg.StoragePath.User,
		cfg.StoragePath.Password,
		cfg.StoragePath.Host,
		cfg.StoragePath.DBName)
	return
}
