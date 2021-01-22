// https://lnly.hatenablog.com/entry/2020/02/26/225722
package router

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/taketake25/IMOFTH/pkg/handler"
)

// "database/sql"

type Page struct {
	Title string
}

func Build() *httprouter.Router {
	router := httprouter.New()

	router.GET("/createImage", handler.createImage)
	router.NotFound = http.HandlerFunc(handler.ApiNotFound)
	router.MethodNotAllowed = http.HandlerFunc(handler.ApiNotFound)

	return router
}
