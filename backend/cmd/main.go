package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"buma/internal/httphandler"
	"buma/internal/image"
	"buma/internal/postgres"
	"buma/internal/s3repo"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(err)
	}

	ctx := context.Background()

	client, err := s3repo.NewS3Client(ctx)
	if err != nil {
		log.Fatal(err)
	}

	pgUrl := os.Getenv("PG_URL")
	db, err := postgres.NewPostgresPool(ctx, pgUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	s3rep := s3repo.NewS3Repository(client)

	imageR := postgres.NewImageRepo(db)
	imageS := image.NewImageService(imageR, s3rep)
	imageH := httphandler.NewImageHandler(imageS)

	r := http.NewServeMux()

	r.HandleFunc("POST /upload", imageH.Upload)
	r.HandleFunc("GET /my/images", imageH.GetImages)
	r.HandleFunc("DELETE /images/{name}", imageH.DeleteImage)

	log.Fatal(http.ListenAndServe(":8080", r))
}
