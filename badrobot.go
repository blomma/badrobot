package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/blomma/badrobot/models"
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
	t, err := template.New("badfriends").Parse(tpl)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	badfriends, err := models.GetAllBadFriends()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	jsonBadFriends, err := json.Marshal(badfriends)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	p := Page{BadFriends: template.JS(jsonBadFriends)}
	t.Execute(w, p)
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

func main() {
	http.Handle("/badfriends",
		gziphandler.GzipHandler(http.HandlerFunc(logHandler(BadFriendsHandler))))

	srv := &http.Server{
		Addr:         ":8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
