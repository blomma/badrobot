package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

func LogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		x, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}

		log.Println(fmt.Sprintf("%q", x))
		defer log.Println("<------")
		next.ServeHTTP(w, r)
	})
}
