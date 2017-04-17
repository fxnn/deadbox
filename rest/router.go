package rest

import (
	"encoding/json"
	"fmt"
	"github.com/fxnn/deadbox/model"
	"github.com/gorilla/mux"
	"net/http"
)

type router struct {
	http.Handler
	drop model.Drop
}

func newRouter(drop model.Drop) *router {
	handler := mux.NewRouter()
	result := &router{Handler: handler, drop: drop}

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

	return result
}

func (r *router) workerId(rq *http.Request) model.WorkerId {
	return model.WorkerId(mux.Vars(rq)["workerId"])
}
func (r *router) outputJson(rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", "application/json")
}
func (r *router) requestInvalid(rw http.ResponseWriter, err error) {
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(rw, "Your request was invalid: %s", err)
}
func (r *router) internalServerError(rw http.ResponseWriter, err error) {
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(rw, "Couldn't handle your request: %s", err)
}

func (r *router) handleGetAllWorkers(
	rw http.ResponseWriter,
	rq *http.Request,
) {
	r.outputJson(rw)
	result, err := r.drop.Workers()
	if err != nil {
		r.internalServerError(rw, fmt.Errorf("couldn't get workers: %s", err))
		return
	}

	err = json.NewEncoder(rw).Encode(result)
	if err != nil {
		r.internalServerError(rw, fmt.Errorf("couldn't serialize workers: %s", err))
		return
	}
}

func (r *router) handlePutWorker(
	rw http.ResponseWriter,
	rq *http.Request,
) {
	var worker *model.Worker = &model.Worker{}
	if err := json.NewDecoder(rq.Body).Decode(worker); err != nil {
		r.requestInvalid(rw, err)
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
	r.outputJson(rw)
	var workerId model.WorkerId = r.workerId(rq)
	result := r.drop.WorkerRequests(workerId)
	json.NewEncoder(rw).Encode(result)
}

func (r *router) handlePutWorkerRequest(
	rw http.ResponseWriter,
	rq *http.Request,
) {
	var request *model.WorkerRequest
	json.NewDecoder(rq.Body).Decode(request)
	r.drop.PutWorkerRequest(request)
}
