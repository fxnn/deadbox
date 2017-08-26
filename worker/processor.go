package worker

import (
	"fmt"
	"log"

	"github.com/fxnn/deadbox/model"
)

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
		// @todo #13 process requests. Create an architecture that can have flexible processors
		log.Printf("received request %s", req.Id)
	}

	return nil
}
