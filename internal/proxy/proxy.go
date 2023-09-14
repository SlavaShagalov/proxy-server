package proxy

import (
	"log"
	"net/http"
)

type Proxy struct {
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Method: %s\n", r.Method)
	log.Printf("RequestURI: %s\n", r.RequestURI)
	log.Printf("Proto: %s\n", r.Proto)
	log.Printf("ContentLength: %d\n", r.ContentLength)
	log.Printf("Header: %s\n", r.Header)
	log.Printf("Host: %s\n", r.Host)
	log.Printf("RemoteAddr: %s\n", r.RemoteAddr)
	log.Printf("URL: %s\n", r.URL)

	r.Header.Del("Proxy-Connection")

}
