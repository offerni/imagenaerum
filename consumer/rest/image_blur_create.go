package rest

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func (srv *Server) ImageBlurCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("hello!!")

	_ = r.MultipartForm.File["files"]
	sigma := r.FormValue("sigma")
	_, err := strconv.ParseFloat(sigma, 64)
	if err != nil {
		log.Println(err)
	}

	// Call RabbitMQ Handler here

	w.WriteHeader(http.StatusNoContent)
}
