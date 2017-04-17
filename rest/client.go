package rest

import (
	"encoding/json"
	"fmt"
	"github.com/fxnn/deadbox/model"
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
		return nil, fmt.Errorf("got response during GET request to \"%s\": %s", address, err)
	}

	var workers []model.Worker
	err = json.NewDecoder(resp.Body).Decode(workers)
	if err != nil {
		return nil, fmt.Errorf("couldn't decode server response: %s", err)
	}

	return workers, nil
}

func (*client) PutWorker(*model.Worker) error {
	panic("implement me")
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
