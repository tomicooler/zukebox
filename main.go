// main.go
package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/tomicooler/zukebox/routes"
)

func main() {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)

	router.Mount("/player/tracks", routes.RoutesTracks())
	router.Mount("/player/recent-tracks", routes.RoutesRecentTracks())
	router.Mount("/player/control", routes.RoutesControl())

	defer routes.PlayerCtx.Release()

	log.Fatal(http.ListenAndServe(":5000", router))
}
