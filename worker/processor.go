package worker

import (
	"fmt"
	"log"

	"encoding/json"

	"github.com/fxnn/deadbox/model"
)

const contentTypePlainText = "text/plain"
const contentTypeJson = "application/json"

type processor struct {
	id   model.WorkerId
	drop model.Drop
}

func (p *processor) pollRequests() error {
	requests, err := p.drop.WorkerRequests(p.id)
	if err != nil {
		return fmt.Errorf("drop returned error: %s", err)
	}

	for _, req := range requests {
		// @todo #7 never process a request twice
		log.Printf("received request %s", req.Id)
		if err = p.enqueueRequest(req); err != nil {
			p.sendErrorResponse(req, err)
		}
	}

	return nil
}

func (p *processor) enqueueRequest(req model.WorkerRequest) error {
	if req.ContentType != contentTypeJson {
		return fmt.Errorf("ContentType not understood by this worker: %s", req.ContentType)
	}

	var q *queueItem = &queueItem{}
	if err := json.Unmarshal(req.Content, q); err != nil {
		return fmt.Errorf("content could not be unmarshalled: %s", err)
	}

	p.addQueueItem(q)
	return nil
}

func (p *processor) addQueueItem(q *queueItem) {
	// @todo #7 create a queue of items to process and process them
}

func (p *processor) sendErrorResponse(r model.WorkerRequest, errToSend error) {
	resp := &model.WorkerResponse{
		Timeout:     r.Timeout,
		ContentType: contentTypePlainText,
		Content:     []byte(errToSend.Error()),
	}
	if err := p.drop.PutWorkerResponse(p.id, r.Id, resp); err != nil {
		log.Printf("drop didn't accept my error response for request %s: %s", r.Id, err)
	}
}
