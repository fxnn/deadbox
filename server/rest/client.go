package rest

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/fxnn/deadbox/server/crypto"
	"github.com/fxnn/deadbox/server/model"
)

// client implements a REST client to a drop server.
type client struct {
	baseUrl *url.URL
	wrapped *http.Client
}

// NewClient creates a new client for the given base URL.
// When verifyByFingerprint is set, regular certificate validation is disabled for TLS connections.
// Instead, the client expects a certificate (may be self-signed) for a public key with the given fingerprint
// configuration.
func NewClient(url *url.URL, verifyByFingerprint *crypto.VerifyByFingerprint) model.Drop {
	wrapped := &http.Client{}
	if verifyByFingerprint != nil {
		tlsConfig := &tls.Config{
			InsecureSkipVerify:    true,
			VerifyPeerCertificate: verifyByFingerprint.VerifyPeerCertificate,
		}
		wrapped.Transport = &http.Transport{TLSClientConfig: tlsConfig}
	}

	return &client{baseUrl: url, wrapped: wrapped}
}

func (c *client) Workers() (workers []model.Worker, err error) {
	err = c.get("worker", &workers)
	return
}
func (c *client) PutWorker(w *model.Worker) error {
	return c.post("worker", w)
}

func (c *client) WorkerRequests(workerId model.WorkerId) (requests []model.WorkerRequest, err error) {
	path := fmt.Sprintf("worker/%s/request", workerId)
	err = c.get(path, &requests)
	return
}

func (c *client) PutWorkerRequest(workerId model.WorkerId, request *model.WorkerRequest) error {
	path := fmt.Sprintf("worker/%s/request", workerId)
	return c.post(path, request)
}

func (c *client) WorkerResponse(workerId model.WorkerId, requestId model.WorkerRequestId) (response model.WorkerResponse, err error) {
	path := fmt.Sprintf("worker/%s/response/%s", workerId, requestId)
	err = c.get(path, &response)
	return
}

func (c *client) PutWorkerResponse(workerId model.WorkerId, requestId model.WorkerRequestId, response *model.WorkerResponse) error {
	path := fmt.Sprintf("worker/%s/response/%s", workerId, requestId)
	return c.post(path, response)
}

func (c *client) post(path string, source interface{}) error {
	var (
		err error
		v   []byte
	)

	if v, err = json.Marshal(source); err != nil {
		return fmt.Errorf("POST request to '%s' couldn't be encoded: %s", path, err)
	}

	address := c.resolveAddress(path)
	resp, err := c.wrapped.Post(address, "application/json", bytes.NewReader(v))
	if err != nil {
		return fmt.Errorf("POST request to '%s' failed: %s", path, err)
	}

	return c.assertStatusOk(resp)
}

func (c *client) assertStatusOk(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		var bodyStr = "<response body not available>"
		if bodyBytes, err := ioutil.ReadAll(resp.Body); err == nil {
			bodyStr = string(bodyBytes)
		}
		return fmt.Errorf("%s request to '%s' returned code %d: %s", resp.Request.Method, resp.Request.URL, resp.StatusCode, bodyStr)
	}

	return nil
}

func (c *client) get(path string, target interface{}) error {
	var err error

	address := c.resolveAddress(path)
	resp, err := c.wrapped.Get(address)
	if err != nil {
		return fmt.Errorf("GET request to '%s' failed: %s", address, err)
	}

	if err := c.assertStatusOk(resp); err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		return fmt.Errorf("GET response from '%s' couldn't be decoded: %s", address, err)
	}

	return nil
}

func (c *client) resolveAddress(path string) string {
	return c.resolveUrl(path).String()
}
func (c *client) resolveUrl(path string) *url.URL {
	referenceUrl, err := url.Parse(path)
	if err != nil {
		panic(fmt.Errorf("couldn't parse \"%s\" as URL: %s", path, err))
	}

	return c.baseUrl.ResolveReference(referenceUrl)
}
