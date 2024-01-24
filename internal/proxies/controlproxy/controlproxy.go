package controlproxy

import (
	"dash/pkg/https"
	"dash/pkg/networksimulator"
	"dash/pkg/noisegenerator"
	"dash/pkg/proxy"
	"dash/pkg/types"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type ControlProxy struct {
	proxy.Proxy
	errChan    chan error
	HTTPServer https.HTTP
}

func (cp *ControlProxy) GetRoutes() []types.RouteInfo {
	return []types.RouteInfo{
		{
			HandlerFunc: cp.Forward,
			Path:        "/{delayPath:.*}",
			Description: "forward request, without modification, to DASH server",
		},
	}
}

func (cp *ControlProxy) Forward(w http.ResponseWriter, r *http.Request) {
	destination := r.URL.Path
	if strings.Contains(r.URL.Path, "delay") {
		destination = fmt.Sprintf(
			"http://dashserver:8080/%s",
			strings.Join(strings.Split(r.URL.Path, "/")[2:], "/"))
	}

	resp, err := cp.HTTPServer.GenericMethod(destination)
	if err != nil {
		cp.errChan <- err
	}

	if strings.Contains(destination, "hlsmanifest") {
		// Simulate network conditions
		networksimulator.NetworkDelay()

		// Add noise to data to reduce signature.
		noisyWriter := noisegenerator.New(w)
		noisyWriter.Write(resp.Bytes)
		log.Println("*******************")
		log.Println("*******************")
		log.Println("*******************")
		log.Println("*******************")

		//w.Write(resp.Bytes)
		//if _, err := io.Copy(w, bytes.NewReader(resp.Bytes)); err != nil {
		//	cp.errChan <- err
		//	http.Error(w, "Error streaming response", http.StatusInternalServerError)
		//	return
		//}

	} else {
		w.Write(resp.Bytes)
	}

}

func New(errChan chan error) ControlProxy {
	proxy := ControlProxy{
		errChan: errChan,
		HTTPServer: https.HTTP{
			Method: "GET",
		},
	}
	return proxy
}
