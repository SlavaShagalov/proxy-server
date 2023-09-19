package proxy

import (
	"github.com/SlavaShagalov/proxy-server/internal/requests"
	"go.uber.org/zap"
	"net/http"
)

type Proxy struct {
	client *http.Client
	rep    requests.Repository
	log    *zap.Logger
}

func New(rep requests.Repository, log *zap.Logger) *Proxy {
	return &Proxy{
		client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		rep: rep,
		log: log,
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.log.Debug("NEW",
		zap.String("method", r.Method),
		zap.String("host", r.Host),
		zap.String("remote_addr", r.RemoteAddr))

	if r.Method == http.MethodConnect {
		p.httpsHandle(w, r)
	} else {
		p.httpHandle(w, r)
	}
}
