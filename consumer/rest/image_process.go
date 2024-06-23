package rest

import (
	"net/http"

	"github.com/offerni/imagenaerum/consumer/img"
)

func (srv *Server) ImageProcess(w http.ResponseWriter, r *http.Request) {
	opts := RequestOpts{
		Request: r,
	}

	if err := opts.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err := srv.ImgSvc.Process(img.ProcessOpts{
		Files:  r.MultipartForm.File["files"],
		Params: r.FormValue("sigma"),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusNoContent)
}

type RequestOpts struct {
	*http.Request
}

func (opts *RequestOpts) Validate() error {
	if err := opts.ParseMultipartForm(10 << 20); err != nil {
		return ErrInvalidForm
	}

	if opts.MultipartForm == nil || (len(opts.MultipartForm.File) == 0 && len(opts.MultipartForm.Value) == 0) {
		return ErrEmptyForm
	}

	return nil
}
