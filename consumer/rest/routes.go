package rest

import (
	"github.com/go-chi/chi"
)

func initializeRoutes(mux *chi.Mux, srv Server) {
	mux.Post("/process", srv.ImageProcess)
	mux.Get("/files/{id}", srv.FileFetch)
}
