package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/Financial-Times/go-fthealth/v1a"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

type alphavilleSeriesHandler struct {
	service alphavilleSeriesService
}

// HealthCheck does something
func (h *alphavilleSeriesHandler) HealthCheck() v1a.Check {
	return v1a.Check{
		BusinessImpact:   "Unable to respond to request for the alphavilleSeries data from TME",
		Name:             "Check connectivity to TME",
		PanicGuide:       "https://sites.google.com/a/ft.com/ft-technology-service-transition/home/run-book-library/alphavilleSeries-transfomer",
		Severity:         1,
		TechnicalSummary: "Cannot connect to TME to be able to supply alphavilleSeries",
		Checker:          h.checker,
	}
}

// Checker does more stuff
func (h *alphavilleSeriesHandler) checker() (string, error) {
	err := h.service.checkConnectivity()
	if err == nil {
		return "Connectivity to TME is ok", err
	}
	return "Error connecting to TME", err
}

func newAlphavilleSeriesHandler(service alphavilleSeriesService) alphavilleSeriesHandler {
	return alphavilleSeriesHandler{service: service}
}

func (h *alphavilleSeriesHandler) getAlphavilleSeries(writer http.ResponseWriter, req *http.Request) {
	obj, found := h.service.getAlphavilleSeries()
	writeJSONResponse(obj, found, writer)
}

func (h *alphavilleSeriesHandler) getAlphavilleSeriesIds(writer http.ResponseWriter, req *http.Request) {
	obj, _ := h.service.getAlphavilleSeriesIds()
	streamJSONResponse(obj, writer)
}

func streamJSONResponse(ids []idEntry, writer http.ResponseWriter) {
	if len(ids) < 1 {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	for _, k := range ids {
		enc := json.NewEncoder(writer)
		enc.Encode(k)
	}

}

func (h *alphavilleSeriesHandler) getAlphavilleSeriesCount(writer http.ResponseWriter, req *http.Request) {
	count := h.service.getAlphavilleSeriesCount()
	fmt.Fprintf(writer, "%d", count)
}

func (h *alphavilleSeriesHandler) getAlphavilleSeriesByUUID(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	uuid := vars["uuid"]

	obj, found := h.service.getAlphavilleSeriesByUUID(uuid)
	writeJSONResponse(obj, found, writer)
}

//GoodToGo returns a 503 if the healthcheck fails - suitable for use from varnish to check availability of a node
func (h *alphavilleSeriesHandler) GoodToGo(writer http.ResponseWriter, req *http.Request) {
	if _, err := h.checker(); err != nil {
		writer.WriteHeader(http.StatusServiceUnavailable)
	}
}

func writeJSONResponse(obj interface{}, found bool, writer http.ResponseWriter) {
	writer.Header().Add("Content-Type", "application/json")

	if !found {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	enc := json.NewEncoder(writer)
	if err := enc.Encode(obj); err != nil {
		log.Errorf("Error on json encoding=%v\n", err)
		writeJSONError(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func writeJSONError(w http.ResponseWriter, errorMsg string, statusCode int) {
	w.WriteHeader(statusCode)
	fmt.Fprintln(w, fmt.Sprintf("{\"message\": \"%s\"}", errorMsg))
}
