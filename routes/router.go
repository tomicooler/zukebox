// router.go
package routes

import (
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func RoutesTracks() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/{trackID}", GetATrack)
	router.Delete("/{trackID}", DeleteTrack)
	router.Group(func(r chi.Router) {
		r.Use(middleware.ThrottleBacklog(5, 50, 60*time.Second))
		r.Post("/", CreateTrack)
	})
	router.Get("/", GetAllTrack)
	return router
}

func RoutesRecentTracks() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/{trackID}", GetARecentTrack)
	router.Get("/", GetAllRecentTrack)
	return router
}

func RoutesControl() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", GetControl)
	router.Patch("/", PatchControl)
	return router
}
