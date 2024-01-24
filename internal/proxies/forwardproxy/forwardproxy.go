package forwardproxy

import (
	"bufio"
	"dash/pkg/https"
	"dash/pkg/proxy"
	"dash/pkg/types"
	"fmt"
	"net/http"
	"os"
)

type ForwardProxy struct {
	proxy.Proxy
	errChan    chan error
	HTTPServer https.HTTP
}

func (cp *ForwardProxy) GetRoutes() []types.RouteInfo {
	return []types.RouteInfo{
		{
			HandlerFunc: cp.Process,
			Path:        "/process",
			Description: "process the most recent trace",
		},
	}
}

func readNginxAccessLog(logFilePath string) error {
	file, err := os.Open(logFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (cp *ForwardProxy) Process(w http.ResponseWriter, r *http.Request) {
	//	READ NGINX ACCESS LOG.
	readNginxAccessLog("/var/log/nginx/nginx.stdout.log")
}

func New(errChan chan error) ForwardProxy {
	proxy := ForwardProxy{
		errChan: errChan,
		HTTPServer: https.HTTP{
			Method: "GET",
		},
	}
	return proxy
}
