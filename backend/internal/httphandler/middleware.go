package httphandler

import (
	"net/http"
	"time"

	"buma/internal/cache"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body = b
	return rw.ResponseWriter.Write(b)
}

func Idempotency(cache *cache.RedisCache, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Idempotency-Key")
		if key == "" {
			next.ServeHTTP(w, r)
			return
		}

		data, err := cache.Get(r.Context(), key)
		if err != nil {
			httpError(w, err)
			return
		}

		if data != nil {
			w.Write(data)
			return
		}

		ok, err := cache.SetNx(r.Context(), key, []byte(`1`), time.Minute*3)
		if err != nil {
			httpError(w, err)
			return
		}

		if !ok {
			httpError(w, errIdempotency)
			return
		}

		response := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(response, r)

		err = cache.Set(r.Context(), key, response.body, time.Hour*24)
		if err != nil {
			httpError(w, err)
			return
		}
	})
}
