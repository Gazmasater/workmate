package health

import (
	"net/http"

	"github.com/gaz358/myprog/workmate/pkg/logger"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	log := logger.Global().Named("health")
	log.Debugw("health check request", "method", r.Method, "path", r.URL.Path, "remote", r.RemoteAddr)

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("ok")); err != nil {
		log.Warnw("failed to write health response", "err", err)
	}
}
