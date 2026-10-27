package httphandler

import (
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

	err = h.service.Upload(r.Context(), "test123", handler.Filename, data)
	if err != nil {
		httpError(w, err)
		return
	}
}

func (h *imageHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	err := h.service.DeleteImage(r.Context(), "test123", name)
	if err != nil {
		httpError(w, err)
		return
	}
}
