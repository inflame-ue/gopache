package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/inflame-ue/gopache/internal/cache"
	"github.com/inflame-ue/gopache/internal/config"
	"github.com/inflame-ue/gopache/internal/proxy"
)

func main() {
	proxyConfig, err := config.NewProxyConfig()
	if err != nil {
		log.Fatal(err)
	}

	cacheMap := cache.NewCacheMap()
	client := &http.Client{}
	if proxyConfig.FlushCache {
		cacheMap.Flush()
		log.Print("cache flushed succesfully")
	}

	proxy := proxy.NewProxy(client, &cacheMap, proxyConfig.Origin)
	addr := fmt.Sprintf(":%d", proxyConfig.Port)
	log.Printf("listening on port: %d", proxyConfig.Port)
	http.ListenAndServe(addr, proxy)
}
