package config

import (
	"errors"
	"flag"
	"net/url"
)

type ProxyConfig struct {
	Port       int
	Origin     string
	FlushCache bool
	CachePath  string
}

func isValidUrl(origin string) bool {
	parsed, err := url.Parse(origin)
	return err == nil && parsed.Scheme != "" && parsed.Host != ""
}

func NewProxyConfigFromFlags(port int, origin string, clearCache bool, cachePath string) (*ProxyConfig, error) {
	proxyConfig := ProxyConfig{
		Port:       port,
		Origin:     origin,
		FlushCache: clearCache,
		CachePath:  cachePath,
	}

	if proxyConfig.FlushCache {
		return &proxyConfig, nil
	}

	if len(origin) == 0 {
		return nil, errors.New("err: missing origin, while --clear-cache is false")
	}

	if !isValidUrl(origin) {
		return nil, errors.New("err: invalid URL format")
	}

	return &proxyConfig, nil
}

func NewProxyConfig() (*ProxyConfig, error) {
	port := flag.Int("port", 8080, "the port on which the proxy server will run (default: 8080)")
	origin := flag.String("origin", "", "the URL of the server to which the request will be forwarded")
	clearCache := flag.Bool("clear-cache", false, "clear the proxy cache, will force all request to be forwarded to origin")
	cachePath := flag.String("cache-path", "cache.json", "path to the cache file, where the persistent cache will be stored (default: \"cache.json\")")
	flag.Parse()

	proxyConfig, err := NewProxyConfigFromFlags(*port, *origin, *clearCache, *cachePath)
	if err != nil {
		return nil, err
	}

	return proxyConfig, nil
}
