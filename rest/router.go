package rest

import (
	"encoding/json"
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

func (r *router) handleGetAllWorkers(
	rw http.ResponseWriter,
	rq *http.Request,
) {
	r.outputJson(rw)
	result := r.drop.Workers()
	json.NewEncoder(rw).Encode(asJsonWorkers(result))
}

func (r *router) handlePutWorker(
	rw http.ResponseWriter,
	rq *http.Request,
) {
	var worker *jsonWorker
	json.NewDecoder(rq.Body).Decode(worker)
	r.drop.PutWorker(worker)
}

func (r *router) handleGetAllWorkerRequests(
	rw http.ResponseWriter,
	rq *http.Request,
) {
	r.outputJson(rw)
	var workerId model.WorkerId = r.workerId(rq)
	result := r.drop.WorkerRequests(workerId)
	json.NewEncoder(rw).Encode(asJsonWorkerRequests(result))
}

func (r *router) handlePutWorkerRequest(
	rw http.ResponseWriter,
	rq *http.Request,
) {
	var request *jsonWorkerRequest
	json.NewDecoder(rq.Body).Decode(request)
	r.drop.PutWorkerRequest(request)
}
