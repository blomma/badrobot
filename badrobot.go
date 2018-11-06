package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/blomma/badrobot/handlers"
)

// Version is the version number or commit hash
// These variables should be set by the linker when compiling
var (
	Version     = "0.0.0"
	CommitHash  = "Unknown"
	CompileDate = "Unknown"
)

// Options
var (
	flagRedis   = flag.String("redis", "", "redis server to hook into")
	flagVersion = flag.Bool("version", false, "Show the version number and information")
)

func commandLineFlags() {
	flag.Parse()
	if *flagVersion {
		fmt.Println("Version:", Version)
		os.Exit(0)
	}
}

func main() {
	commandLineFlags()

	badFriendsHandler := handlers.NewBadFriendsHandler(flagRedis)
	mux := http.NewServeMux()
	mux.HandleFunc("/badfriends", badFriendsHandler.Handler)

	srv := &http.Server{
		Handler:      gziphandler.GzipHandler(handlers.LogHandler(mux)),
		Addr:         ":8001",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		badFriendsHandler.Stop()
		os.Exit(1)
	}()

	log.Println("Version:", Version)
	log.Println("Commit hash:", CommitHash)
	log.Println("Compiled on:", CompileDate)

	log.Fatal(srv.ListenAndServe())
}
