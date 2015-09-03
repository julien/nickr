package utils

import "net/http"

func AddCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")
		w.Header().Set("Access-Control-Allow-Methods", "GET,DELETE,PATCH,POST,PUT")
		next.ServeHTTP(w, r)
	})
}
