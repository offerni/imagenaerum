package rest

import (
	"net/http"

	"github.com/offerni/imagenaerum/consumer/img"
)

func (srv *Server) ImageProcess(w http.ResponseWriter, r *http.Request) {
	err := srv.ImgSvc.Process(img.ProcessOpts{
		Files:  r.MultipartForm.File["files"],
		Params: r.FormValue("sigma"),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusNoContent)
}
