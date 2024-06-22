package rest

import "github.com/go-chi/chi"

func initializeRoutes(mux *chi.Mux, srv Server) {
	mux.Post("/blur", srv.ImageBlur)
}
