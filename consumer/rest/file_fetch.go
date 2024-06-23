package rest

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/offerni/imagenaerum/worker/utils"
)

func (srv *Server) FileFetch(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	filePath := filepath.Join(utils.ConvertedPath, id)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error checking file", http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, filePath)
}
