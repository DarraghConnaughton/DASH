package proxy

import (
	"dash/pkg/server"
	"dash/pkg/types"
	"log"
	"os/exec"
)

type Proxy struct {
	server.Server
	errChan chan error
}

func (p *Proxy) Start(bind string, routes []types.RouteInfo) {
	if len(routes) > 0 {
		p.Server.LoadRoutes(routes)
		go p.Server.Start(bind)
	}
	log.Println("32423423432423432423423")
	log.Println("32423423432423432423423")
	log.Println("32423423432423432423423")
	p.LaunchNGINX()
}

func (p *Proxy) LaunchNGINX() {
	_, err := exec.Command(
		"nginx",
		"-c",
		"/etc/nginx/nginx.conf").CombinedOutput()

	if err != nil {
		p.errChan <- err
	}
}

func New(errChan chan error) Proxy {
	proxy := Proxy{
		errChan: errChan,
	}
	return proxy
}
