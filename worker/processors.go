package worker

import (
	"github.com/fxnn/deadbox/config"
	"github.com/fxnn/deadbox/request"
	"github.com/fxnn/deadbox/request/echo"
)

type requestProcessors struct {
	processorsById map[string]request.Processor
}

func (r *requestProcessors) requestProcessorForId(requestProcessorId string) (request.Processor, bool) {
	p, ok := r.processorsById[requestProcessorId]
	return p, ok
}

func createRequestProcessorsByIdMap(c *config.Worker) (result map[string]request.Processor) {
	result = make(map[string]request.Processor)
	addRequestProcessor(result, echo.New())
	return
}

func addRequestProcessor(processors map[string]request.Processor, processor request.Processor) {
	processors[processor.Id()] = processor
}
