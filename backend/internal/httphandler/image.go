package httphandler

import (
	"encoding/json"
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
	file, handler, err := r.FormFile("image")
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

	imageType := http.DetectContentType(data)

	err = checkContentType(imageType)
	if err != nil {
		httpError(w, err)
		return
	}

	err = h.service.Upload(r.Context(), "images", handler.Filename, imageType, data)
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
	name := r.PathValue("name")

	err := h.service.DeleteImage(r.Context(), "images", name)
	if err != nil {
		httpError(w, err)
		return
	}
}
