package worker

import (
	"fmt"
)

type RequestProcessor interface {
	RequestType() string
	EmptyRequestContent() interface{}
	ProcessRequest(requestContent interface{}) interface{}
}

type requestProcessors struct {
	processorsByRequestType map[string]RequestProcessor
}

func (r *requestProcessors) AddRequestProcessor(p RequestProcessor) error {
	t := p.RequestType()
	if _, ok := r.processorsByRequestType[t]; ok {
		return fmt.Errorf("there is already a processor for RequestType: %s", t)
	}

	r.processorsByRequestType[t] = p
	return nil
}

func (r *requestProcessors) requestProcessorForType(requestType string) (RequestProcessor, bool) {
	p, ok := r.processorsByRequestType[requestType]
	return p, ok
}
