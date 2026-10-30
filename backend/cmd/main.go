package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"buma/internal/cache"
	"buma/internal/httphandler"
	"buma/internal/image"
	"buma/internal/postgres"
	"buma/internal/s3repo"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
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

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:       0,
		Protocol: 2,
	})
	defer redisClient.Close()

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

	r.HandleFunc("POST /api/upload", imageH.Upload)
	r.HandleFunc("GET /api/my/images", imageH.GetImages)
	r.HandleFunc("DELETE /api/images/{y}/{m}/{name}", imageH.DeleteImage)

	r.HandleFunc("POST /test", func(w http.ResponseWriter, r *http.Request) {
		sleep, err := strconv.Atoi(r.URL.Query().Get("s"))
		if err != nil {
			sleep = 2
		}

		time.Sleep(time.Duration(sleep) * time.Second)

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"status": "ok"}`))
	})

	server := &http.Server{
		Addr:    ":8081",
		Handler: httphandler.Idempotency(cache.NewRedisCache(redisClient), r),
	}

	log.Fatal(server.ListenAndServeTLS("cert.crt", "cert.pem"))
}
