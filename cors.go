package main

import "net/http"

// AddCORS adds Cross Origin Resource Sharing headers to responses
func AddCORS(next http.Handler, origin, headers, methods string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if len(origin) == 0 {
			origin = "*"
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)

		if len(headers) > 0 {
			w.Header().Set("Access-Control-Allow-Headers", headers)
		}

		if len(methods) > 0 {
			w.Header().Set("Access-Control-Allow-Methods", methods)
		}
		next.ServeHTTP(w, r)
	})
}
