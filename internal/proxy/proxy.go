package proxy

import (
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

	for name, hdrs := range resp.Header {
		for _, hdr := range hdrs {
			w.Header().Add(name, hdr)
		}
	}
	w.WriteHeader(resp.StatusCode)

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("error while copying the response body: %v", err)
		return
	}
}
