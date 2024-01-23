package proxy

import (
	"dash/pkg/https"
	"dash/pkg/server"
	"dash/pkg/types"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

type ProxyServer struct {
	server.Server
	errChan    chan error
	HTTPServer https.HTTPS
}

func (p *ProxyServer) getRoutes() []types.RouteInfo {
	return []types.RouteInfo{
		{
			HandlerFunc: p.Forward,
			//Path:        "/delay/{delayPath:.*}",
			Path:        "/{delayPath:.*}",
			Description: "forward request, without modification, to DASH server",
		},
	}
}

func (p *ProxyServer) Delay() {

	fmt.Println("blahblah")
}

func (p *ProxyServer) Obfuscate() {

	fmt.Println("blahblah")
}

func (p *ProxyServer) Start(bind string) {
	log.Println("we are here.")

	log.Println(p.getRoutes())
	log.Println(bind)
	log.Println("p.getRoutes()p.getRoutes()p.getRoutes()")
	//p.Server.LoadRoutes(p.getRoutes())
	p.Server.LoadRoutes(p.getRoutes())
	p.Server.Start(bind)
}

func (p *ProxyServer) Forward(w http.ResponseWriter, r *http.Request) {

	log.Println("234234234234234234")
	log.Println("234234234234234234")
	log.Println("234234234234234234")
	log.Println(r.URL)
	log.Println(r.URL.Path)
	destination := r.URL.Path
	if strings.Contains(r.URL.Path, "delay") {
		destination = fmt.Sprintf("http://dashserver:8080/%s", strings.Join(strings.Split(r.URL.Path, "/")[2:], "/"))
	}

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate a random number between 0 and 999
	randomNumber := rand.Intn(1000)

	// Bias towards the upper end (e.g., over 500, 70% of the time)
	if rand.Float64() < 0.7 {
		randomNumber += 500
	}

	time.Sleep(time.Duration(randomNumber) * time.Millisecond)
	resp, err := p.HTTPServer.GenericMethod(destination)
	fmt.Println("resprespresprespresp")
	fmt.Println(resp)
	fmt.Println(err)
	w.Write(resp.Bytes)
}

func (p *ProxyServer) LaunchNGINX() {
	cmd := exec.Command("nginx", "-c", "/etc/nginx/nginx.conf")

	// Run the command and capture its output
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(output))
		log.Println("Error executing command:", err)
		p.errChan <- err
	}

	// Print the output of the command
	log.Println("Command Output:", string(output))
}

func New(errChan chan error) ProxyServer {
	proxy := ProxyServer{
		errChan: errChan,
		HTTPServer: https.HTTPS{
			Method: "GET",
		},
	}
	return proxy
}
