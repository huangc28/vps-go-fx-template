package main

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	cachefx "vps-go-fx-template/cache/fx"
	dbfx "vps-go-fx-template/db/fx"
	appfx "vps-go-fx-template/internal/app/fx"
	healthfx "vps-go-fx-template/internal/app/health/fx"
	routerfx "vps-go-fx-template/internal/router/fx"
	serverfx "vps-go-fx-template/internal/server/fx"
)

func main() {
	app := fx.New(
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),
		appfx.CoreAppOptions,
		dbfx.Module,
		cachefx.Module,
		routerfx.CoreRouterOptions,
		serverfx.ServerOptions,
		healthfx.Module,
	)

	app.Run()
}
