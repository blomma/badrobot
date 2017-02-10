package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/blomma/badrobot/models"
	"github.com/pkg/profile"
)

var (
	g_dir = flag.String("dir", ".", "the directory to serve files from. Defaults to the current dir")
)

const tpl = `
<!DOCTYPE html>
<html>
	<head>
	<meta charset="UTF-8">
	<title>{{.Title}}</title>
	</head>
	<body>
	{{range .Items}}<div>{{ .  }}</div>{{else}}<div><strong>no rows</strong></div>{{end}}
	</body>
</html>`

func BadFriendsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	// t, err := template.New("badfriends").Parse(tpl)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Load json data

	//
}

func main() {
	defer profile.Start().Stop()
	flag.Parse()

	badfriends, err := models.GetAllBadFriends()
	if err != nil {
		log.Fatal(err)
	}

	for _, value := range badfriends {
		log.Print(value)
	}
	// r := mux.NewRouter()
	// s := r.Host("badfriends.artsoftheinsane.com")
	// // This will serve files under http://localhost:8000/static/<filename>
	// // r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(g_dir))))
	// r.HandleFunc("/badfriends", BadFriendsHandler)

	// srv := &http.Server{
	// 	Handler: r,
	// 	Addr:    "127.0.0.1:8000",
	// 	// Good practice: enforce timeouts for servers you create!
	// 	WriteTimeout: 15 * time.Second,
	// 	ReadTimeout:  15 * time.Second,
	// }

	// log.Fatal(srv.ListenAndServe())
}
