package model

import "time"

// Drop defines the interface for a drop, able to manage workers and store their
// requests and responses.
type Drop interface {
	Workers() []Worker
	PutWorker(*Worker)

	WorkerRequests(WorkerId) []WorkerRequest
	PutWorkerRequest(*WorkerRequest)

	WorkerResponse(WorkerRequestId) []WorkerResponse
	PutWorkerResponse(*WorkerResponse)
}

type WorkerId string
type Worker struct {
	Id      WorkerId
	Timeout time.Time
}

type WorkerRequestId string
type WorkerRequest struct {
	Id      WorkerRequestId
	Timeout time.Time
}
type WorkerResponse struct {
	Id      WorkerRequestId
	Timeout time.Time
}
