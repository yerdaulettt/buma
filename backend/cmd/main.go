package main

import (
	"context"
	"log"
	"net/http"

	"buma/internal/httphandler"
	"buma/internal/image"
	"buma/internal/s3repo"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	client, err := s3repo.NewS3Client(ctx)
	if err != nil {
		log.Fatal(err)
	}

	s3rep := s3repo.NewS3Repository(client)

	imageS := image.NewImageService(nil, s3rep)
	imageH := httphandler.NewImageHandler(imageS)

	r := http.NewServeMux()

	r.HandleFunc("POST /upload", imageH.Upload)
	r.HandleFunc("DELETE /images/{name}", imageH.DeleteImage)

	log.Fatal(http.ListenAndServe(":8080", r))
}
