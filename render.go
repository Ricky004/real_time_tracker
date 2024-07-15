package main

import (
	"net/http"
	"html/template"

)

func renderTemplate(w http.ResponseWriter, tmpl string) {
	t, err := template.ParseFiles("view/" + tmpl + ".html")
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    err = t.Execute(w, nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func handler(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, "index")
}