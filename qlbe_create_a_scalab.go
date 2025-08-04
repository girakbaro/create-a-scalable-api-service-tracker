package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// ServiceTracker represents a service tracker
type ServiceTracker struct {
	services  map[string]int
	serviceMu sync.RWMutex
}

// NewServiceTracker returns a new service tracker
func NewServiceTracker() *ServiceTracker {
	return &ServiceTracker{
		services:  make(map[string]int),
		serviceMu: sync.RWMutex{},
	}
}

// TrackService increments the count for a given service
func (st *ServiceTracker) TrackService(service string) {
	st.serviceMu.Lock()
	defer st.serviceMu.Unlock()
	st.services[service]++
}

// GetServiceCount returns the count for a given service
func (st *ServiceTracker) GetServiceCount(service string) int {
	st.serviceMu.RLock()
	defer st.serviceMu.RUnlock()
	return st.services[service]
}

// API represents the API service tracker
type API struct {
	tracker *ServiceTracker
}

// NewAPI returns a new API service
func NewAPI() *API {
	return &API{
		tracker: NewServiceTracker(),
	}
}

// TrackServiceHandler handles tracking a service
func (a *API) TrackServiceHandler(w http.ResponseWriter, r *http.Request) {
	service := mux.Vars(r)["service"]
	a.tracker.TrackService(service)
	w.WriteHeader(http.StatusNoContent)
}

// GetServiceCountHandler handles getting the service count
func (a *API) GetServiceCountHandler(w http.ResponseWriter, r *http.Request) {
	service := mux.Vars(r)["service"]
	count := a.tracker.GetServiceCount(service)
	json, err := json.Marshal(struct {
		Count int `json:"count"`
	}{
		count,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func main() {
	api := NewAPI()
	r := mux.NewRouter()
	r.HandleFunc("/track/{service}", api.TrackServiceHandler).Methods("POST")
	r.HandleFunc("/count/{service}", api.GetServiceCountHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))
}