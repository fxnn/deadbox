package rest

import (
	"crypto/tls"
	"net/http"
)

// TLS abstracts the way REST interfaces TLS is configured.
type TLS interface {
	// Config provides the tls.Config object to be set on server.TLSConfig
	Config() (*tls.Config, error)
	// ListenAndServe encapsulates the decision, whether server.ListenAndServe() or server.ListenAndServeTLS() should
	// be invoked.
	ListenAndServe(*http.Server) error
}

type noTLS struct{}

func NoTLS() TLS {
	return &noTLS{}
}
func (t *noTLS) Config() (*tls.Config, error) {
	return nil, nil
}
func (t *noTLS) ListenAndServe(server *http.Server) error {
	return server.ListenAndServe()
}

type fileBasedTLS struct {
	keyFile  string
	certFile string
}

func NewFileBasedTLS(keyFile string, certFile string) TLS {
	return &fileBasedTLS{keyFile: keyFile, certFile: certFile}
}
func (t *fileBasedTLS) Config() (config *tls.Config, err error) {
	config = &tls.Config{
		Certificates: make([]tls.Certificate, 1),
		NextProtos:   []string{"h2", "http/1.1"}, // HINT: enable HTTP/2
	}

	config.Certificates[0], err = tls.LoadX509KeyPair(t.certFile, t.keyFile)
	return
}
func (t *fileBasedTLS) ListenAndServe(server *http.Server) error {
	return server.ListenAndServeTLS("", "")
}
