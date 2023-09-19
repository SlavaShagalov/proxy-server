package proxy

import (
	"log"
	"net/http"
	"net/url"
	"strings"
)

func addrToUrl(addr string) *url.URL {
	if !strings.HasPrefix(addr, "https") {
		addr = "https://" + addr
	}
	u, err := url.Parse(addr)
	if err != nil {
		log.Fatal(err)
	}
	return u
}

func changeRequestToTarget(req *http.Request, targetHost string) {
	targetUrl := addrToUrl(targetHost)
	targetUrl.Path = req.URL.Path
	targetUrl.RawQuery = req.URL.RawQuery
	req.URL = targetUrl
	req.RequestURI = ""
}

func parseCookies(cookies []*http.Cookie) map[string]string {
	cookieData := make(map[string]string)
	for _, cookie := range cookies {
		cookieData[cookie.Name] = cookie.Value
	}
	return cookieData
}
