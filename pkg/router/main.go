// https://lnly.hatenablog.com/entry/2020/02/26/225722
package router

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	// "github.com/taketake25/IMOFTH/pkg/handler"
)

// "database/sql"

func Build() {
	router := httprouter.New()
	router.GET("/createImage", handler.createImage)
	router.NotFound = http.FileServer(http.Dir("public"))

	return router
}
