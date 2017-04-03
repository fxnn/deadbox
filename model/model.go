package model

import "time"

// Drop defines the interface for a drop, able to manage workers and store their
// requests and responses.
type Drop interface {
	Workers() []Worker
	PutWorker(Worker)

	WorkerRequests(WorkerId) []WorkerRequest
	PutWorkerRequest(WorkerRequest)

	WorkerResponse(WorkerRequestId) []WorkerResponse
	PutWorkerResponse(WorkerResponse)
}

type WorkerId string
type Worker interface {
	Id() WorkerId
	Timeout() time.Time
}

type WorkerRequestId string
type WorkerRequest interface {
	Id() WorkerRequestId
	Timeout() time.Time
}
type WorkerResponse interface {
	Id() WorkerRequestId
	Timeout() time.Time
}
