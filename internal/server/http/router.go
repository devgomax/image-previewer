package internalhttp

import (
	"net/http"

	"github.com/devgomax/image-previewer/internal/app/image_previewer"
	"github.com/go-chi/chi/v5"
)

// NewRouter creates a new HTTP router with the specified middleware.
func NewRouter(app *imagepreviewer.App, middlewares ...func(http.Handler) http.Handler) http.Handler {
	r := chi.NewRouter()
	for _, m := range middlewares {
		r.Use(m)
	}

	r.Get("/fill/{width}/{height}/*", app.PreviewImage)

	return r
}
