package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func newProxy(serviceUrl string) (*httputil.ReverseProxy, error) {
	sUrl, err := url.Parse(serviceUrl)
	if err != nil {
		return nil, err
	}

	rp := httputil.NewSingleHostReverseProxy(sUrl)
	rp.Director = func(r *http.Request) {
		r.URL.Scheme = sUrl.Scheme
		r.URL.Host = sUrl.Host
		r.Host = sUrl.Host

		r.Header.Set("X-Forwarded-For", r.RemoteAddr)
	}

	return rp, nil
}

func GetRoutes(services map[string][]string) (*http.ServeMux, error) {
	r := http.NewServeMux()

	for sUrl, endpoints := range services {
		rp, err := newProxy(sUrl)
		if err != nil {
			return nil, err
		}

		for _, e := range endpoints {
			r.HandleFunc(e, rp.ServeHTTP)
		}
	}

	return r, nil
}
