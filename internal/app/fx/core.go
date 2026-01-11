package fx

import (
	"vps-go-fx-template/config"
	"vps-go-fx-template/internal/logs"

	"go.uber.org/fx"
)

var CoreAppOptions = fx.Options(
	fx.Provide(
		config.NewViper,
		config.NewConfig,
		logs.NewLogger,
		logs.NewSugaredLogger,
	),
)
