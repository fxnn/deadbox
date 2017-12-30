package model

import (
	"time"
)

// Drop defines the interface for a drop, able to manage workers and store their
// requests and responses.
type Drop interface {
	// Workers is invoked by a user, so that he knows which workers he can interact with.
	Workers() ([]Worker, error)
	// PutWorker is invoked by a worker to register himself at the drop.
	PutWorker(*Worker) error

	// WorkerRequests is invoked by a worker to retrieve the list of all open requests he shall process.
	WorkerRequests(WorkerId) ([]WorkerRequest, error)
	// PutWorkerRequest is invoked by a user to transfer a request to a worker.
	PutWorkerRequest(WorkerId, *WorkerRequest) error

	// WorkerResponse is invoked by a user to receive the worker's response on a request, and to know that the request
	// has been processed successfully.
	WorkerResponse(WorkerId, WorkerRequestId) (WorkerResponse, error)
	// PutWorkerResponse is invoked by a worker to transfer a response to a successfully processed request.
	PutWorkerResponse(WorkerId, WorkerRequestId, *WorkerResponse) error
}

type WorkerId string
type Worker struct {
	Id      WorkerId
	Name    string
	Timeout time.Time

	// PublicKey contains the public RSA key of the worker, marshalled in ASN1 format.
	PublicKey []byte
}

type WorkerRequestId string
type WorkerRequest struct {
	Id      WorkerRequestId
	Timeout time.Time
	Content []byte

	// ContentType describes the kind of data contained in Content.
	ContentType string

	// EncryptionType describes the kind of encryption applied to Content.
	EncryptionType string
}
type WorkerResponse struct {
	Timeout     time.Time
	Content     []byte
	ContentType string
}
