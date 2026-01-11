package fx

import (
	"vps-go-fx-template/internal/app/health"
	"vps-go-fx-template/internal/router"

	"go.uber.org/fx"
)

var Module = fx.Options(
	router.AsRoute(health.NewHandler),
)
