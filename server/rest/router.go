package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fxnn/deadbox/server/model"
	"github.com/gorilla/mux"
)

type router struct {
	handler http.Handler
	drop    model.Drop
}

func newRouter(drop model.Drop) *router {
	handler := mux.NewRouter()
	result := &router{handler: handler, drop: drop}

	handler.Path("/worker").
		HandlerFunc(result.handleGetAllWorkers).
		Methods("GET")
	handler.Path("/worker").
		HandlerFunc(result.handlePutWorker).
		Methods("POST")
	handler.Path("/worker/{workerId}/request").
		HandlerFunc(result.handleGetAllWorkerRequests).
		Methods("GET")
	handler.Path("/worker/{workerId}/request").
		HandlerFunc(result.handlePutWorkerRequest).
		Methods("POST")
	handler.Path("/worker/{workerId}/response/{requestId}").
		HandlerFunc(result.handleGetWorkerResponse).
		Methods("GET")
	handler.Path("/worker/{workerId}/response/{requestId}").
		HandlerFunc(result.handlePutWorkerResponse).
		Methods("POST")

	return result
}

func (r *router) ServeHTTP(w http.ResponseWriter, r2 *http.Request) {
	r.handler.ServeHTTP(w, r2)
}

func (r *router) workerId(rq *http.Request) (model.WorkerId, error) {
	workerId := mux.Vars(rq)["workerId"]
	if workerId == "" {
		return "", fmt.Errorf("workerId must be set")
	}
	return model.WorkerId(workerId), nil
}
func (r *router) requestId(rq *http.Request) (model.WorkerRequestId, error) {
	requestId := mux.Vars(rq)["requestId"]
	if requestId == "" {
		return "", fmt.Errorf("requestId must be set")
	}
	return model.WorkerRequestId(requestId), nil
}
func (r *router) outputJson(rw http.ResponseWriter, v interface{}) error {
	rw.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(v)
}
func (r *router) inputJson(rq *http.Request, v interface{}) error {
	return json.NewDecoder(rq.Body).Decode(v)
}

func (r *router) requestInvalid(rw http.ResponseWriter, err error) {
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(rw, "your request was invalid: %s", err)
}
func (r *router) internalServerError(rw http.ResponseWriter, err error) {
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(rw, "couldn't handle your request: %s", err)
}

func (r *router) handleGetAllWorkers(
	rw http.ResponseWriter,
	rq *http.Request,
) {
	result, err := r.drop.Workers()
	if err != nil {
		r.internalServerError(rw, fmt.Errorf("couldn't get workers: %s", err))
		return
	}

	err = r.outputJson(rw, result)
	if err != nil {
		r.internalServerError(rw, fmt.Errorf("couldn't serialize workers: %s", err))
		return
	}
}

func (r *router) handlePutWorker(
	rw http.ResponseWriter,
	rq *http.Request,
) {
	worker := &model.Worker{}
	if err := r.inputJson(rq, worker); err != nil {
		r.requestInvalid(rw, fmt.Errorf("couldn't read worker: %s", err))
		return
	}
	if err := r.drop.PutWorker(worker); err != nil {
		r.internalServerError(rw, fmt.Errorf("couldn't put worker: %s", err))
		return
	}
}

func (r *router) handleGetAllWorkerRequests(
	rw http.ResponseWriter,
	rq *http.Request,
) {
	workerId, err := r.workerId(rq)
	if err != nil {
		r.requestInvalid(rw, err)
	}

	result, err := r.drop.WorkerRequests(workerId)
	if err != nil {
		r.internalServerError(rw, fmt.Errorf("couldn't get worker requests: %s", err))
		return
	}

	if err := r.outputJson(rw, result); err != nil {
		r.internalServerError(rw, fmt.Errorf("couldn't encode worker request: %s", err))
	}
}

func (r *router) handlePutWorkerRequest(
	rw http.ResponseWriter,
	rq *http.Request,
) {
	var (
		workerId model.WorkerId
		request  *model.WorkerRequest
		err      error
	)

	if workerId, err = r.workerId(rq); err != nil {
		r.requestInvalid(rw, fmt.Errorf("no workerId given"))
		return
	}

	request = &model.WorkerRequest{}
	if err = r.inputJson(rq, request); err != nil {
		r.requestInvalid(rw, fmt.Errorf("couldn't read request: %s", err))
		return
	}

	if err := r.drop.PutWorkerRequest(workerId, request); err != nil {
		r.internalServerError(rw, fmt.Errorf("couldn't put worker request: %s", err))
	}
}

func (r *router) handleGetWorkerResponse(
	rw http.ResponseWriter,
	rq *http.Request,
) {
	workerId, err := r.workerId(rq)
	if err != nil {
		r.requestInvalid(rw, err)
	}
	requestId, err := r.requestId(rq)
	if err != nil {
		r.requestInvalid(rw, err)
		return
	}

	result, err := r.drop.WorkerResponse(workerId, requestId)
	if err != nil {
		r.internalServerError(rw, fmt.Errorf("couldn't get worker response: %s", err))
		return
	}

	if err := r.outputJson(rw, result); err != nil {
		r.internalServerError(rw, fmt.Errorf("couldn't encode worker response: %s", err))
	}
}

func (r *router) handlePutWorkerResponse(
	rw http.ResponseWriter,
	rq *http.Request,
) {
	var (
		workerId  model.WorkerId
		requestId model.WorkerRequestId
		response  *model.WorkerResponse
		err       error
	)

	if workerId, err = r.workerId(rq); err != nil {
		r.requestInvalid(rw, err)
		return
	}
	if requestId, err = r.requestId(rq); err != nil {
		r.requestInvalid(rw, err)
		return
	}

	response = &model.WorkerResponse{}
	if err = r.inputJson(rq, response); err != nil {
		r.requestInvalid(rw, fmt.Errorf("couldn't read response: %s", err))
		return
	}

	if err := r.drop.PutWorkerResponse(workerId, requestId, response); err != nil {
		r.internalServerError(rw, fmt.Errorf("couldn't put worker response: %s", err))
	}
}
