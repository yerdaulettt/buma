package httphandler

import (
	"net/http"
)

func httpError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"error": "` + err.Error() + `"}`))
}
