package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/Novochenko/sso/internal/app/grpc"
	"github.com/Novochenko/sso/internal/services/auth"
	"github.com/Novochenko/sso/internal/storage/mysql"
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

	storage, err := mysql.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)
	return &App{
		GRPCServer: grpcApp,
	}
}
