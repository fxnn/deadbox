package processor

type RequestProcessor interface {
	// RequestType returns the string by which this RequestProcessor is identified.
	RequestType() string
	// EmptyRequestContent creates an empty instance of accepted request content. It is passed into the unmarshalling
	// framework, and then given to the ProcessRequest function.
	EmptyRequestContent() interface{}
	// ProcessRequest processes the given request. The result is passed into the marshalling framework.
	ProcessRequest(requestContent interface{}) interface{}
}
