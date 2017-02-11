package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/blomma/badrobot/models"
	"github.com/gorilla/mux"
)

var (
	g_dir = flag.String("dir", ".", "the directory to serve files from. Defaults to the current dir")
)

const tpl = `
<!DOCTYPE html>
<html>
<head>
<style>
#map {
	height: 100%;
	width: 100%;
}
html,
body {
	height: 100%;
	margin: 0;
	padding: 0;
}
</style>
</head>
<body>
<div id="map"></div>
<script>
function initMap() {
	var badFriendsData = {{.BadFriends}}
	var map = new google.maps.Map(document.getElementById('map'), {
		zoom: 1,
		center: new google.maps.LatLng(2.8, -187.3),
		mapTypeId: 'terrain'
	});

	var markers = badFriendsData.map(function(location, i) {
		return new google.maps.Marker({
			position: new google.maps.LatLng(location.latitude, location.longitude),
			map: map,
			icon: {
                path: google.maps.SymbolPath.CIRCLE,
                scale: 5,
                fillColor: 'red',
                fillOpacity: .2,
                strokeColor: 'white',
                strokeWeight: .5
            }
        });
    });
}
</script>
<script async defer
src="https://maps.googleapis.com/maps/api/js?key=AIzaSyCdqrB2bNdayZDaNNJqkUKTmzTH4DUtmco&callback=initMap">
</script>
</body>
</html>`

type Page struct {
	BadFriends template.JS
}

func BadFriendsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	t, err := template.New("badfriends").Parse(tpl)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	badfriends, err := models.GetAllBadFriends()
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	jsonBadFriends, err := json.Marshal(badfriends)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	p := Page{BadFriends: template.JS(jsonBadFriends)}

	t.Execute(w, p)
}

func main() {
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/badfriends", BadFriendsHandler)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
