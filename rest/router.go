package rest

import (
	jsonenc "encoding/json"
	"fmt"
	"github.com/fxnn/deadbox/json"
	"github.com/fxnn/deadbox/model"
	"github.com/gorilla/mux"
	"log"
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
	rw.WriteHeader(400)
	fmt.Fprintf(rw, "Your request was invalid: %v", err)
}

func (r *router) handleGetAllWorkers(
	rw http.ResponseWriter,
	rq *http.Request,
) {
	r.outputJson(rw)
	result := r.drop.Workers()
	err := jsonenc.NewEncoder(rw).Encode(json.AsWorkers(result))
	if err != nil {
		log.Println("Couldn't serialize worker:", err)
	}
}

func (r *router) handlePutWorker(
	rw http.ResponseWriter,
	rq *http.Request,
) {
	var worker *json.Worker = &json.Worker{}
	if err := jsonenc.NewDecoder(rq.Body).Decode(worker); err != nil {
		r.requestInvalid(rw, err)
		return
	}
	r.drop.PutWorker(worker)
}

func (r *router) handleGetAllWorkerRequests(
	rw http.ResponseWriter,
	rq *http.Request,
) {
	r.outputJson(rw)
	var workerId model.WorkerId = r.workerId(rq)
	result := r.drop.WorkerRequests(workerId)
	jsonenc.NewEncoder(rw).Encode(json.AsWorkerRequests(result))
}

func (r *router) handlePutWorkerRequest(
	rw http.ResponseWriter,
	rq *http.Request,
) {
	var request *json.WorkerRequest
	jsonenc.NewDecoder(rq.Body).Decode(request)
	r.drop.PutWorkerRequest(request)
}
