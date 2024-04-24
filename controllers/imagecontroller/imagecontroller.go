package imagecontroller

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	http.ServeFile(w, r, "uploads/"+filename)
}
