package main

import (
	"log"
	"net/http"

	"api-gateway/internal/proxy"
)

func main() {
	services := map[string][]string{
		"http://localhost:8333":  {"GET /images/"},
		"https://localhost:8081": {"GET /api/my/images", "POST /api/upload", "DELETE /api/images/{y}/{m}/{name}"},
	}

	r, err := proxy.GetRoutes(services)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Start")
	// log.Fatal(http.ListenAndServe(":8080", proxy.LogMiddleware(r)))
	log.Fatal(http.ListenAndServeTLS(":8080", "cert.crt", "cert.pem", proxy.LogMiddleware(r)))
}
