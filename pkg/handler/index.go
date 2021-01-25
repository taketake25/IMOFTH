// https://qiita.com/taizo/items/bf1ec35a65ad5f608d45
package handler

import (
	"net/http"
	"text/template"
)

type Page struct {
	Title string
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Del("Content-Type")
	w.WriteHeader(200)

	page := Page{"Hello"}
	tmpl, err := template.ParseFiles("../index.html") // ParseFilesを使う
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, page)
	if err != nil {
		panic(err)
	}
}
