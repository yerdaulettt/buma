package httphandler

import (
	"errors"
	"net/http"
	"slices"
)

var (
	errContentType = errors.New("Incorrect content type")
)

func httpError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"error": "` + err.Error() + `"}`))
}

func checkContentType(cType string) error {
	correctTypes := []string{"image/jpeg", "image/png"}

	if !slices.Contains(correctTypes, cType) {
		return errContentType
	}

	return nil
}
