package forwardproxy

import (
	"bufio"
	"dash/pkg/https"
	"dash/pkg/proxy"
	"dash/pkg/types"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type ForwardProxy struct {
	proxy.Proxy
	errChan    chan error
	HTTPServer https.HTTP
	logRegex   *regexp.Regexp
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

func (cp *ForwardProxy) readNginxAccessLog(logFilePath string) []types.NetworkTraceData {
	var networkTraces []types.NetworkTraceData
	file, err := os.Open(logFilePath)
	if err != nil {
		cp.errChan <- err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		matches := cp.logRegex.FindStringSubmatch(scanner.Text())
		fmt.Println("matchesmatchesmatchesmatchesmatches")
		fmt.Println(matches)
		if len(matches) == 4 {

			log.Println("matches")
			log.Println(matches)
			log.Println("matches")

			pathParts := strings.Split(strings.Split(matches[2], "HTTP")[0], "/")

			networkTraces = append(networkTraces, types.NetworkTraceData{
				Timestamp:  matches[1],
				Resolution: pathParts[len(pathParts)-2],
				Sequence:   pathParts[len(pathParts)-3],
				Bytes:      matches[3],
			})

			log.Println("*************")
			log.Println(matches[1])
			log.Println(pathParts[len(pathParts)-2])
			log.Println(pathParts[len(pathParts)-3])
			log.Println(matches[3])
			log.Println("*************")
		}
	}
	if err := scanner.Err(); err != nil {
		cp.errChan <- err
	}
	return networkTraces
}

func (cp *ForwardProxy) Process(w http.ResponseWriter, r *http.Request) {
	//	READ NGINX ACCESS LOG.
	log.Println("reretretretretrtreterte")
	log.Println("reretretretretrtreterte")
	log.Println("reretretretretrtreterte")
	log.Println("reretretretretrtreterte")
	log.Println("reretretretretrtreterte")
	log.Println("reretretretretrtreterte")
	log.Println("reretretretretrtreterte")
	tmp := cp.readNginxAccessLog("/var/log/nginx/nginx.stdout.log")
	data, err := json.Marshal(tmp)
	if err != nil {
		cp.errChan <- err
	}
	w.Write(data)

}

func New(errChan chan error) ForwardProxy {
	proxy := ForwardProxy{
		errChan: errChan,
		HTTPServer: https.HTTP{
			Method: "GET",
		},
		logRegex: regexp.MustCompile(`\[([\d.]+)\] "(GET [^"]+)" \d+ (\d+)`),
	}
	return proxy
}
