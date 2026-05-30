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
}

func isValidUrl(origin string) bool {
	parsed, err := url.Parse(origin)
	return err == nil && parsed.Scheme != "" && parsed.Host != ""
}

func NewProxyConfig() (*ProxyConfig, error) {
	port := flag.Int("port", 8080, "the port on which the proxy server will run (default: 8080)")
	origin := flag.String("origin", "", "the URL of the server to which the request will be forwarded")
	clearCache := flag.Bool("clear-cache", false, "clear the proxy cache, will force all request to be forwarded to origin")
	flag.Parse()

	proxyConfig := ProxyConfig{
		Port:       *port,
		Origin:     *origin,
		FlushCache: *clearCache,
	}
	if proxyConfig.FlushCache {
		return &proxyConfig, nil
	}

	if len(*origin) == 0 {
		return &ProxyConfig{}, errors.New("err: missing origin, while --clear-cache is false")
	}
	
	if !isValidUrl(*origin) {
		return &ProxyConfig{}, errors.New("err: invalid URL format")
	}

	return &proxyConfig, nil
}
