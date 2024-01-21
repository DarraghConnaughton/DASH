package https

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type HTTPS struct {
	TLSConfig TLSConfig   `json:"tls_config"`
	Header    http.Header `json:"header"`
	Method    string      `json:"method"`
	Body      io.Reader   `json:"body"`
}

type TLSConfig struct {
	CAFile   string
	CertFile string
	KeyFile  string
}

func (config *TLSConfig) GetTLSConfig() (*tls.Config, error) {
	// Load CA certificate
	caCert, err := ioutil.ReadFile(config.CAFile)
	if err != nil {
		return nil, err
	}
	// Load certificate and private key
	cert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
	if err != nil {
		return &tls.Config{}, err
	}
	// Create a certificate pool and add the CA certificate
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	return &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
	}, nil
}

func (h *HTTPS) GenericMethod(hostname string) (HTTPResponse, error) {
	serverResponse := newResponse()
	// Prepare Request
	req, err := http.NewRequest(h.Method, hostname, h.Body)
	if err != nil {
		return serverResponse, err
	}
	req.Header = h.Header
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}

	//*************
	// To be enabled.
	//config, err := h.TLSConfig.GetTLSConfig()
	//if err != nil {}
	//dashclient := &http.Client{Transport: &http.Transport{
	//	TLSClientConfig: config,
	//}}
	//to be enabed
	//*************
	startTime := time.Now()
	if resp, err := client.Do(req); err != nil {
		serverResponse.RTT = time.Now().Sub(startTime)
		return serverResponse, err
	} else {
		serverResponse.extractRawBytes(resp)
		serverResponse.extractResponseDetails(resp)
		return serverResponse, nil
	}
}
