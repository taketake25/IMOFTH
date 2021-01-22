package handler

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func ApiNotFound(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Del("Content-Type")
	w.WriteHeader(404)

	return
}
