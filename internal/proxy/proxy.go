package proxy

import (
	"bytes"
	"fmt"
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

func addHeaders(w http.ResponseWriter, headers http.Header) {
	for name, hdrs := range headers {
		for _, hdr := range hdrs {
			w.Header().Add(name, hdr)
		}
	}
}

func writeCachedResponse(w http.ResponseWriter, cachedResponse cache.CachedResponse) error {
	w.Header().Add("X-Cache", "HIT")
	addHeaders(w, cachedResponse.Headers)
	w.WriteHeader(cachedResponse.Status)

	_, err := io.Copy(w, bytes.NewReader(cachedResponse.Body))
	if err != nil {
		return fmt.Errorf("error while copying the response body: %v", err)
	}
	return nil
}

func (p *Proxy) performRequest(w http.ResponseWriter, r *http.Request, url string) (*http.Response, []byte, error) {
	redirectReq, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		return &http.Response{}, []byte{}, fmt.Errorf("error while creating the redirect request: %v", err)
	}

	resp, err := p.Client.Do(redirectReq)
	if err != nil {
		return &http.Response{}, []byte{}, fmt.Errorf("error while performing the request to origin: %v", err)
	}
	defer resp.Body.Close()

	w.Header().Add("X-Cache", "MISS")
	addHeaders(w, resp.Header)
	w.WriteHeader(resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &http.Response{}, []byte{}, fmt.Errorf("error while reading the response body bytes into memory: %v", err)
	}

	_, err = io.Copy(w, bytes.NewReader(body))
	if err != nil {
		return &http.Response{}, []byte{}, fmt.Errorf("error while copying the response body: %v", err)
	}

	return resp, body, nil
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	redirectURL, err := url.JoinPath(p.Origin, r.URL.Path)
	if err != nil {
		log.Printf("error while constructing the routed url path: %v", err)
		return
	}

	if cachedResp, ok := p.Cache.Get(redirectURL); ok {
		log.Printf("cache hit for %v", redirectURL)
		err := writeCachedResponse(w, *cachedResp)
		if err != nil {
			log.Print(err)
		}
		return
	}

	log.Printf("redirecting to: %v", redirectURL)
	resp, bodyBytes, err := p.performRequest(w, r, redirectURL)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("handled the request with status: %v", resp.StatusCode)

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		log.Printf("setting a cache entry for: %v", redirectURL)

		cachedResp := cache.CachedResponse{
			Status:  resp.StatusCode,
			Headers: resp.Header,
			Body:    bodyBytes,
		}
		p.Cache.Set(redirectURL, &cachedResp)
	}
}
