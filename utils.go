package nickr

import (
	"log"
	"net/http"
)

// Debugging is a bool type used to enable or disable logging.
type Debugging bool

// Printf logs a message given a format, it uses the standard "log" package.
func (d Debugging) Printf(format string, argv ...interface{}) {
	if d {
		log.Printf(format, argv...)
	}
}

// AddCORS adds Cross Origin Resource Sharing headers to responses.
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
