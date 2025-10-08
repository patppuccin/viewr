package server

import (
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/patppuccin/viewr/src/helpers"
	"github.com/patppuccin/viewr/src/include"
	"github.com/patppuccin/viewr/src/models"
)

func setupRoutes(serverCtx *models.AppContext) (http.Handler, error) {

	// Initialize router
	r := chi.NewRouter()

	// Add custom middlewares
	r.Use(loadServerContext(serverCtx))
	r.Use(logRequest())

	// Add chi middlewares
	r.Use(middleware.CleanPath)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Compress(5,
		"text/html",
		"text/css",
		"application/javascript",
		"application/json",
		"application/wasm",
		"application/xml",
		"text/plain",
		"text/javascript",
		"image/svg+xml",
	))

	// Mount Assets as a Static Route (using embedded assets)
	assetsRootFS, err := fs.Sub(include.Assets, "assets")
	if err != nil {
		return r, helpers.SafeErr("failed to locate embedded assets", err)
	}

	assetsServer := http.FileServer(http.FS(assetsRootFS))

	assetsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Basic cache headers
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Header().Add("Vary", "Accept-Encoding")

		// Serve static asset
		assetsServer.ServeHTTP(w, r)
	})

	// Strip "/assets/" prefix for proper path resolution
	r.Handle("/assets/*", http.StripPrefix("/assets/", assetsHandler))

	// TODO: Mount Page & Media Route Handlers

	// Return router
	return r, nil
}
