package app

import (
	"log"
	"net/http"

	"github.com/taketake25/IMOFTH/pkg/router"
)

func init() {
	r := router.Build()
	log.Fatal(http.ListenAndServe(":8080", r))
}
