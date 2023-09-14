package main

import (
	"github.com/SlavaShagalov/proxy-server/internal/proxy"
	"log"
	"net/http"
)

const (
	host         = "127.0.0.1"
	port         = "8080"
	proxyAddress = host + ":" + port
)

func main() {
	handler := new(proxy.Proxy)

	log.Printf("Starting proxy server on %s", proxyAddress)
	http.ListenAndServe(proxyAddress, handler)
}
