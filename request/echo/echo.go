package echo

import (
	"github.com/fxnn/deadbox/request"
)

const interfaceVersion = "1.0"
const RequestProcessorId = "github.com/fxnn/deadbox/request/echo " + interfaceVersion

type requestProcessor struct{}

func New() request.Processor {
	return &requestProcessor{}
}

// Id returns the string by which this Processor is identified.
func (p *requestProcessor) Id() string {
	return RequestProcessorId
}

// EmptyContent creates an empty instance of accepted request content.
func (p *requestProcessor) EmptyContent() interface{} {
	return make(map[string]interface{})
}

// Process processes the given request.
func (p *requestProcessor) Process(requestContent interface{}) interface{} {
	return requestContent
}
