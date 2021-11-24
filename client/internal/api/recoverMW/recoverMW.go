package recoverMW

import (
	"fmt"
	"net/http"
)

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("Error panic: %s (%T)\n", err, err)
				}
			}()
			next.ServeHTTP(w, r)
		},
	)
}