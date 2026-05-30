package proxy

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/inflame-ue/gopache/internal/cache"
)

type Proxy struct {
	Origin string
	Client *http.Client
	Cache  cache.Cache
}

func NewProxy(client *http.Client, cache cache.Cache, origin string) *Proxy {
	return &Proxy{
		Client: client,
		Cache:  cache,
		Origin: origin,
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	redirectURL, err := url.JoinPath(p.Origin, r.URL.Path)
	if err != nil {
		log.Printf("error while constructing the routed url path: %v", err)
		return
	}
	log.Printf("redirecting to: %v", redirectURL)

	if cachedResp, ok := p.Cache.Get(redirectURL); ok {
		log.Printf("cache hit for %v", redirectURL)

		w.Header().Add("X-Cache", "HIT")
		for name, hdrs := range cachedResp.Headers {
			for _, hdr := range hdrs {
				w.Header().Add(name, hdr)
			}
		}
		w.WriteHeader(cachedResp.Status)

		_, err = io.Copy(w, bytes.NewReader(cachedResp.Body))
		if err != nil {
			log.Printf("error while copying the response body: %v", err)
			return
		}
		return
	}

	redirectReq, err := http.NewRequest(r.Method, redirectURL, r.Body)
	if err != nil {
		log.Printf("error while creating the redirect request: %v", err)
		return
	}

	resp, err := p.Client.Do(redirectReq)
	if err != nil {
		log.Printf("error while performing the request to origin: %v", err)
		return
	}
	defer resp.Body.Close()

	w.Header().Add("X-Cache", "MISS")
	for name, hdrs := range resp.Header {
		for _, hdr := range hdrs {
			w.Header().Add(name, hdr)
		}
	}
	log.Printf("handled the request with status: %v", resp.StatusCode)
	w.WriteHeader(resp.StatusCode)

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("error while copying the response body: %v", err)
		return
	}

	log.Printf("setting a cache entry for: %v", redirectURL)
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error while reading the response body bytes into memory: %v", err)
		return
	}

	cachedResp := cache.CachedResponse{
		Status:  resp.StatusCode,
		Headers: resp.Header,
		Body:    bodyBytes,
	}
	p.Cache.Set(redirectURL, &cachedResp)
}
