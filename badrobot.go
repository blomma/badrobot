package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/blomma/badrobot/models"
)

type Page struct {
	BadFriends template.JS
}

var g_template *template.Template

func BadFriendsHandler(w http.ResponseWriter, r *http.Request) {
	p := Page{BadFriends: template.JS(models.BadFriends.Get())}
	g_template.Execute(w, p)
}

func logHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		x, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}

		log.Println(fmt.Sprintf("%q", x))
		fn(w, r)
	}
}

func init() {
	filename := "badfriends.html"
	g_template = template.Must(template.ParseFiles(filename))
}

func main() {
	http.Handle("/badfriends",
		gziphandler.GzipHandler(
			http.HandlerFunc(logHandler(BadFriendsHandler))))

	srv := &http.Server{
		Addr:         ":8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
