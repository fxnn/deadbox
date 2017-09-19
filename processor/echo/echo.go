package echo

import (
	"github.com/fxnn/deadbox/processor"
)

const interfaceVersion = "1.0"
const RequestType = "github.com/fxnn/deadbox/processor/echo " + interfaceVersion

type echoRequestProcessor struct{}

func New() processor.RequestProcessor {
	return &echoRequestProcessor{}
}

// RequestType returns the string by which this RequestProcessor is identified.
func (p *echoRequestProcessor) RequestType() string {
	return RequestType
}

// EmptyRequestContent creates an empty instance of accepted request content.
func (p *echoRequestProcessor) EmptyRequestContent() interface{} {
	return make(map[string]interface{})
}

// ProcessRequest processes the given request.
func (p *echoRequestProcessor) ProcessRequest(requestContent interface{}) interface{} {
	return requestContent
}
