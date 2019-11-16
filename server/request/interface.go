package request

type Processor interface {
	// Id returns the string by which this Processor is identified.
	Id() string
	// EmptyContent creates an empty instance of accepted request content. It is passed into the unmarshalling
	// framework, and then given to the Process function.
	EmptyContent() interface{}
	// Process processes the given request. The result is passed into the marshalling framework.
	Process(requestContent interface{}) interface{}
}
