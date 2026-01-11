package fx

import (
	"net/http"
	"time"

	"vps-go-fx-template/internal/router"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var CoreRouterOptions = fx.Options(
	fx.Provide(NewMux),
)

type muxParams struct {
	fx.In

	Logger   *zap.Logger
	Handlers []router.Handler `group:"routes"`
}

func NewMux(p muxParams) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(zapRequestLogger(p.Logger))

	for _, h := range p.Handlers {
		h.RegisterRoute(r)
	}

	return r
}

func zapRequestLogger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()

			next.ServeHTTP(ww, r)

			logger.Info(
				"http_request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", ww.Status()),
				zap.Int("bytes", ww.BytesWritten()),
				zap.Duration("duration", time.Since(start)),
			)
		})
	}
}
