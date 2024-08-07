package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/Novochenko/sso/internal/app/grpc"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	grpcApp := grpcapp.New(log, grpcPort)
	return &App{
		GRPCServer: grpcApp,
	}
}
