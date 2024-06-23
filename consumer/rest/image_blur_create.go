package rest

import (
	"log"
	"net/http"
)

func (srv *Server) ImageBlurCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// files := r.MultipartForm.File["files"]
	// sigma := r.FormValue("sigma")
	// sigmaFoat, err := strconv.ParseFloat(sigma, 64)
	// if err != nil {
	// 	log.Println(err)
	// }

	// if err := srv.ImgService.Blur(img.BlurOpts{
	// 	Files: files,
	// 	Sigma: sigmaFoat,
	// }); err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// }

	w.WriteHeader(http.StatusNoContent)
}
