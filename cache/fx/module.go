package fx

import (
	"vps-go-fx-template/cache"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"redis",
	fx.Provide(cache.NewRedis),
)
