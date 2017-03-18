package main

import (
	"flag"
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

var (
	// Version is the version number or commit hash
	// These variables should be set by the linker when compiling
	Version     = "0.0.0-unknown"
	CommitHash  = "Unknown"
	CompileDate = "Unknown"
)

var (
	templateBadFriends *template.Template
	badFriendsModel    *models.BadFriends
)

func badFriendsHandler(w http.ResponseWriter, r *http.Request) {
	p := pageBadFriends{BadFriends: template.JS(badFriendsModel.Result())}
	templateBadFriends.Execute(w, p)
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Version: %v\n", Version)
	fmt.Fprintf(w, "Commit hash: %v\n", CommitHash)
	fmt.Fprintf(w, "Compiled on: %v\n", CompileDate)
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

var (
	flagVersion     = flag.Bool("version", false, "Show the version number and information")
	flagVersionOnly = flag.Bool("versionOnly", false, "Show only the version number")
)

func commandLineFlags() {
	flag.Parse()
	if *flagVersion {
		fmt.Println("Version:", Version)
		fmt.Println("Commit hash:", CommitHash)
		fmt.Println("Compiled on", CompileDate)
		os.Exit(0)
	}

	if *flagVersionOnly {
		fmt.Println(Version)
		os.Exit(0)
	}
}

func main() {
	commandLineFlags()

	filename := "badfriends.html"
	templateBadFriends = template.Must(template.ParseFiles(filename))
	localBadFriendsModel, stopBadFriends := models.NewBadFriends()
	badFriendsModel = localBadFriendsModel

	mux := http.NewServeMux()
	mux.HandleFunc("/badfriends", badFriendsHandler)
	mux.HandleFunc("/version", versionHandler)

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
		stopBadFriends()
		os.Exit(1)
	}()

	log.Println("Version:", Version)
	log.Println("Commit hash:", CommitHash)
	log.Println("Compiled on:", CompileDate)

	log.Fatal(srv.ListenAndServe())
}
