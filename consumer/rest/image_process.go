package rest

import (
	"net/http"

	"github.com/offerni/imagenaerum/consumer/img"
	"github.com/offerni/imagenaerum/worker/utils"
)

func (srv *Server) ImageProcess(w http.ResponseWriter, r *http.Request) {
	opts := RequestOpts{
		Request: r,
	}

	if err := opts.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var blur *string
	if r.FormValue("blur") != "" {
		blur = utils.ToPointer(r.FormValue("blur"))
	}

	var cropAnchor *string
	if r.FormValue("crop_anchor") != "" {
		cropAnchor = utils.ToPointer(r.FormValue("crop_anchor"))
	}

	var resize *string
	if r.FormValue("resize") != "" {
		resize = utils.ToPointer(r.FormValue("resize"))
	}

	var grayscale *string
	if r.FormValue("grayscale") != "" {
		grayscale = utils.ToPointer(r.FormValue("grayscale"))
	}

	var invert *string
	if r.FormValue("invert") != "" {
		invert = utils.ToPointer(r.FormValue("invert"))
	}

	err := srv.ImgSvc.Process(img.ProcessOpts{
		Blur:       blur,
		CropAnchor: cropAnchor,
		Files:      r.MultipartForm.File["files"],
		Grayscale:  grayscale,
		Invert:     invert,
		Resize:     resize,
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
