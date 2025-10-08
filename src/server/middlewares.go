package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/patppuccin/viewr/src/constants"
	"github.com/patppuccin/viewr/src/models"
)

func loadServerContext(serverCtx *models.AppContext) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), constants.AppCtxKey, serverCtx)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func logRequest() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			serverCtx, ok := r.Context().Value(constants.AppCtxKey).(*models.AppContext)
			if !ok || serverCtx == nil || serverCtx.Logger == nil {
				fmt.Println("Unable to fetch the logger from the server context")
				return
			}
			logger := serverCtx.Logger

			status := ww.Status()
			duration := time.Since(start)

			event := logger.Info()
			switch {
			case status >= 500:
				event = logger.Error()
			case status >= 400:
				event = logger.Warn()
			case status >= 300:
				event = logger.Info()
			case status >= 200:
				event = logger.Debug()
			}

			event.
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Int("status", status).
				Dur("duration", duration).
				Msg("req served")
		})
	}
}
