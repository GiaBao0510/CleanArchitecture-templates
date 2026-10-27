package route

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	mdw "github.com/GiaBao0510/CleanArchitecture-templates/internal/adapter/delivery/rest/middleware"
	user "github.com/GiaBao0510/CleanArchitecture-templates/internal/adapter/delivery/rest/user"
)

func NewRouter(h *user.Handler, logger *mdw.LoggingMiddleware) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(logger.Handler)

	r.Route("api/v1/users", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/", h.List)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})

	return r
}
