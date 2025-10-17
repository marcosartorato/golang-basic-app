package hello

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

// HelloHandler is the handler for the "Hello, World!" path.
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	log := getLogger(r)
	if _, err := fmt.Fprintln(w, "Hello, World!"); err != nil {
		// Log the error; since client likely went away, not much else to do
		log.Error("failed to write response", zap.Error(err))
		return
	}
	log.Debug("handled hello request")
}
