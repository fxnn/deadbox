package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fxnn/deadbox/model"
	"io/ioutil"
	"net/http"
	"net/url"
)

// client implements a REST client to a drop server.
type client struct {
	baseUrl *url.URL
}

func NewClient(url *url.URL) model.Drop {
	return &client{url}
}

func (c *client) resolveAddress(reference string) string {
	return c.resolveUrl(reference).String()
}
func (c *client) resolveUrl(reference string) *url.URL {
	referenceUrl, err := url.Parse(reference)
	if err != nil {
		panic(fmt.Errorf("couldn't parse \"%s\" as URL: %s", reference, err))
	}

	return c.baseUrl.ResolveReference(referenceUrl)
}

func (c *client) Workers() ([]model.Worker, error) {
	var err error

	address := c.resolveAddress("worker")
	resp, err := http.Get(address)
	if err != nil {
		return nil, fmt.Errorf("got error during GET request to \"%s\": %s", address, err)
	}

	var workers []model.Worker
	err = json.NewDecoder(resp.Body).Decode(workers)
	if err != nil {
		return nil, fmt.Errorf("couldn't decode server response: %s", err)
	}

	return workers, nil
}

func (c *client) PutWorker(w *model.Worker) error {
	var err error
	var v []byte

	v, err = json.Marshal(w)
	if err != nil {
		return fmt.Errorf("couldn't encode worker: %s", err)
	}

	address := c.resolveAddress("worker")
	resp, err := http.Post(address, "application/json", bytes.NewReader(v))
	if err != nil {
		return fmt.Errorf("got error during POST request to \"%s\": %s", address, err)
	}

	if resp.StatusCode != http.StatusOK {
		var bodyStr string = "<empty>"
		if bodyBytes, err := ioutil.ReadAll(resp.Body); err == nil {
			bodyStr = string(bodyBytes)
		}
		return fmt.Errorf("got status %d during POST request to \"%s\": %s", resp.StatusCode, address, bodyStr)
	}

	return nil
}

func (*client) WorkerRequests(model.WorkerId) []model.WorkerRequest {
	panic("implement me")
}

func (*client) PutWorkerRequest(*model.WorkerRequest) {
	panic("implement me")
}

func (*client) WorkerResponse(model.WorkerRequestId) []model.WorkerResponse {
	panic("implement me")
}

func (*client) PutWorkerResponse(*model.WorkerResponse) {
	panic("implement me")
}
