package rest

import (
	"net/http"
)

func (srv Server) ImageBlur(w http.ResponseWriter, r *http.Request) {
	if err := srv.ImgService.Blur("./files/raw/imgtest.jpg", 5); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
}
