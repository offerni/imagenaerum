package rest

import (
	"log"
	"net/http"
	"strconv"
)

func (srv Server) ImageBlurCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Println(err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]
	sigma := r.FormValue("sigma")
	sigmaFoat, err := strconv.ParseFloat(sigma, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to parse sigma", http.StatusBadRequest)
	}

	if err := srv.ImgService.Blur(files, sigmaFoat); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusNoContent)
}
