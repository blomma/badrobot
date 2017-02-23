package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/blomma/badrobot/models"
)

type pageBadFriends struct {
	BadFriends template.JS
}

var templateBadFriends *template.Template

func badFriendsHandler(w http.ResponseWriter, r *http.Request) {
	p := pageBadFriends{BadFriends: template.JS(models.BadFriends.Get())}
	templateBadFriends.Execute(w, p)
}

func logHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		x, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}

		log.Println(fmt.Sprintf("%q", x))
		defer log.Println("<------")
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func init() {
	filename := "badfriends.html"
	templateBadFriends = template.Must(template.ParseFiles(filename))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/badfriends", badFriendsHandler)

	srv := &http.Server{
		Handler:      gziphandler.GzipHandler(logHandler(mux)),
		Addr:         ":8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(1)
	}()

	log.Fatal(srv.ListenAndServe())
}
