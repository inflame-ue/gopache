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

	// this is a no-op for now, since persistent cache is not implemented
	if proxyConfig.FlushCache {
		cacheMap.Flush()
		log.Print("cache flushed succesfully")
		os.Exit(0)
	}

	proxy := proxy.NewProxy(client, cacheMap, proxyConfig.Origin)
	addr := fmt.Sprintf(":%d", proxyConfig.Port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Print("interrupt received...serializing cache and exiting...")
		cacheMap.Save("test.json")
		os.Exit(0)
	}()
	
	log.Printf("listening on port: %d", proxyConfig.Port)
	err = http.ListenAndServe(addr, proxy)
	if err != nil {
		log.Print("something went horribly wrong...")
		os.Exit(1)
	}
}
