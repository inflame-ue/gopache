package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

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
	proxy := proxy.NewProxy(client, cacheMap, proxyConfig.Origin)
	addr := fmt.Sprintf(":%d", proxyConfig.Port)

	if proxyConfig.FlushCache {
		cacheMap.Flush(proxyConfig.CachePath)
		log.Print("cache flushed succesfully")
		os.Exit(0)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Print("interrupt received...serializing cache and exiting...")
		cacheMap.Save(proxyConfig.CachePath)
		os.Exit(0)
	}()
	
	log.Printf("listening on port: %d", proxyConfig.Port)
	err = http.ListenAndServe(addr, proxy)
	if err != nil {
		log.Print("something went horribly wrong...")
		os.Exit(1)
	}
}
