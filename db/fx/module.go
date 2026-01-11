package fx

import (
	"vps-go-fx-template/db"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"sqlx-postgres-db",
	fx.Provide(db.NewSQLXPostgresDB),
)
