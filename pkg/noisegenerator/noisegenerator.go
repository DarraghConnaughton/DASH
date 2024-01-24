package noisegenerator

import (
	"math/rand"
	"net/http"
)

type NoisyResponseWriter struct {
	originalWriter http.ResponseWriter
}

func (nrw *NoisyResponseWriter) Write(p []byte) (int, error) {
	noise := []byte("!!thisIsTheDelimiter!!")
	data := make([]byte, rand.Intn(5069044))
	noise = append(noise, data...)
	p = append(p, noise...)
	return nrw.originalWriter.Write(p)
}

func (nrw *NoisyResponseWriter) Header() http.Header {
	return nrw.originalWriter.Header()
}

func (nrw *NoisyResponseWriter) WriteHeader(statusCode int) {
	nrw.originalWriter.WriteHeader(statusCode)
}

func New(w http.ResponseWriter) NoisyResponseWriter {
	return NoisyResponseWriter{
		originalWriter: w,
	}
}
