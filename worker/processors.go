package worker

import (
	"github.com/fxnn/deadbox/config"
	"github.com/fxnn/deadbox/processor"
	"github.com/fxnn/deadbox/processor/echo"
)

type requestProcessors struct {
	processorsByRequestType map[string]processor.RequestProcessor
}

func (r *requestProcessors) requestProcessorForType(requestType string) (processor.RequestProcessor, bool) {
	p, ok := r.processorsByRequestType[requestType]
	return p, ok
}

func createProcessorsByRequestTypeMap(c config.Worker) (result map[string]processor.RequestProcessor) {
	result = make(map[string]processor.RequestProcessor)
	addRequestProcessor(result, echo.New())
	return
}

func addRequestProcessor(processors map[string]processor.RequestProcessor, processor processor.RequestProcessor) {
	processors[processor.RequestType()] = processor
}
