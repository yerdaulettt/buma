package httphandler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"buma/internal/image"
)

type imageHandler struct {
	service image.ImageService
}

func NewImageHandler(s *image.ImageService) *imageHandler {
	return &imageHandler{service: *s}
}

func (h *imageHandler) Upload(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("image")
	if err != nil {
		httpError(w, err)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		httpError(w, err)
		return
	}

	imageType, err := getImageType(http.DetectContentType(data))
	if err != nil {
		httpError(w, err)
		return
	}

	err = h.service.Upload(r.Context(), "images", imageType, data)
	if err != nil {
		httpError(w, err)
		return
	}
}

func (h *imageHandler) GetImages(w http.ResponseWriter, r *http.Request) {
	images, err := h.service.GetImages(r.Context(), 1)
	if err != nil {
		httpError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(&images)
	if err != nil {
		httpError(w, err)
		return
	}
}

func (h *imageHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	year := r.PathValue("y")
	month := r.PathValue("m")
	name := r.PathValue("name")

	err := h.service.DeleteImage(r.Context(), "images", fmt.Sprintf("%s/%s/%s", year, month, name))
	if err != nil {
		httpError(w, err)
		return
	}
}
