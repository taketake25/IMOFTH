// https://qiita.com/taizo/items/bf1ec35a65ad5f608d45
package handler

import (
	"html/template"
	"net/http"
)

type Page struct {
	Title string
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
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
