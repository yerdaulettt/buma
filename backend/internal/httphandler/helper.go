package httphandler

import (
	"errors"
	"net/http"
	"slices"
	"strings"
)

var (
	errContentType = errors.New("Incorrect content type")
)

func httpError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"error": "` + err.Error() + `"}`))
}

func getImageType(contentType string) (string, error) {
	correctTypes := []string{"image/jpeg", "image/png"}

	if !slices.Contains(correctTypes, contentType) {
		return "", errContentType
	}

	imageType := strings.Split(contentType, "/")

	return imageType[1], nil
}
