package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

const testUUID = "bba39990-c78d-3629-ae83-808c333c6dbc"
const oneUUID = "one"
const twoUUID = "two"
const getIdsResponse = `{"id":"` + testUUID + `"}` + "\n"
const getTwoIdsResponse = `{"id":"` + oneUUID + `"}` + "\n" + `{"id":"` + twoUUID + `"}` + "\n"
const countResponse = "1"
const getAlphavilleSeriesResponse = `[{"apiUrl":"http://localhost:8080/transformers/alphavilleseries/bba39990-c78d-3629-ae83-808c333c6dbc"}]`  +  "\n"
const getAlphavilleSeriesByUUIDResponse = `{"uuid":"bba39990-c78d-3629-ae83-808c333c6dbc","alternativeIdentifiers":{"TME":["MTE3-U3ViamVjdHM="],"uuids":["bba39990-c78d-3629-ae83-808c333c6dbc"]},"prefLabel":"Global Alphaville Series","type":"AlphavilleSeries"}` +  "\n"

func TestHandlers(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name         string
		req          *http.Request
		dummyService alphavilleSeriesService
		statusCode   int
		contentType  string // Contents of the Content-Type header
		body         string
	}{
		{"Success - get alphavilleSeries by uuid", newRequest("GET", fmt.Sprintf("/transformers/alphavilleseries/%s", testUUID)), &dummyService{found: true, alphavilleSeries: []alphavilleSeries{getDummyAlphavilleSeries(testUUID, "Global Alphaville Series", "MTE3-U3ViamVjdHM=")}}, http.StatusOK, "application/json", getAlphavilleSeriesByUUIDResponse},
		{"Not found - get alphavilleSeries by uuid", newRequest("GET", fmt.Sprintf("/transformers/alphavilleseries/%s", testUUID)), &dummyService{found: false, alphavilleSeries: []alphavilleSeries{alphavilleSeries{}}}, http.StatusNotFound, "application/json", ""},
		{"Success - get alphavilleSeries", newRequest("GET", "/transformers/alphavilleseries"), &dummyService{found: true, alphavilleSeries: []alphavilleSeries{alphavilleSeries{UUID: testUUID}}}, http.StatusOK, "application/json", getAlphavilleSeriesResponse},
		{"Not found - get alphavilleSeries", newRequest("GET", "/transformers/alphavilleseries"), &dummyService{found: false, alphavilleSeries: []alphavilleSeries{}}, http.StatusNotFound, "application/json", ""},
		{"Success - get alphavilleSeries ids", newRequest("GET", "/transformers/alphavilleseries/__ids"), &dummyService{found: true, alphavilleSeries: []alphavilleSeries{alphavilleSeries{UUID: testUUID}}}, http.StatusOK, "", getIdsResponse},
		{"Success - get 2 alphavilleSeries ids", newRequest("GET", "/transformers/alphavilleseries/__ids"), &dummyService{found: true, alphavilleSeries: []alphavilleSeries{alphavilleSeries{UUID: oneUUID}, alphavilleSeries{UUID: twoUUID}}}, http.StatusOK, "", getTwoIdsResponse},
		{"Not found - get alphavilleSeries", newRequest("GET", "/transformers/alphavilleseries/__ids"), &dummyService{found: false, alphavilleSeries: []alphavilleSeries{}}, http.StatusNotFound, "application/json", ""},
		{"Success - get alphavilleSeries count", newRequest("GET", "/transformers/alphavilleseries/__count"), &dummyService{found: true, alphavilleSeries: []alphavilleSeries{alphavilleSeries{UUID: testUUID}}}, http.StatusOK, "application/json", countResponse},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		router(test.dummyService).ServeHTTP(w, test.req)
		assert.True(test.statusCode == w.Code, fmt.Sprintf("%s: Wrong response code, was %d, should be %d", test.name, w.Code, test.statusCode))
		assert.Equal(test.body, w.Body.String(), fmt.Sprintf("%s: Wrong body", test.name))
	}
}

func newRequest(method, url string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	return req
}

func router(s alphavilleSeriesService) *mux.Router {
	m := mux.NewRouter()
	h := newAlphavilleSeriesHandler(s)
	m.HandleFunc("/transformers/alphavilleseries", h.getAlphavilleSeries).Methods("GET")
	m.HandleFunc("/transformers/alphavilleseries/{uuid:([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})}", h.getAlphavilleSeriesByUUID).Methods("GET")
	m.HandleFunc("/transformers/alphavilleseries/__ids", h.getAlphavilleSeriesIds).Methods("GET")
	m.HandleFunc("/transformers/alphavilleseries/__count", h.getAlphavilleSeriesCount).Methods("GET")
	return m
}

type dummyService struct {
	found            bool
	alphavilleSeries []alphavilleSeries
}

func (s *dummyService) getAlphavilleSeries() ([]alphavilleSeriesLink, bool) {
	var alphavilleSeriesLinks []alphavilleSeriesLink
	for _, sub := range s.alphavilleSeries {
		alphavilleSeriesLinks = append(alphavilleSeriesLinks, alphavilleSeriesLink{APIURL: "http://localhost:8080/transformers/alphavilleseries/" + sub.UUID})
	}
	return alphavilleSeriesLinks, s.found
}

func (s *dummyService) getAlphavilleSeriesByUUID(uuid string) (alphavilleSeries, bool) {
	return s.alphavilleSeries[0], s.found
}

func (s *dummyService) getAlphavilleSeriesCount() int {
	return len(s.alphavilleSeries)
}

func (s *dummyService) getAlphavilleSeriesIds() ([]idEntry, bool) {
	var ids []idEntry
	for _, sub := range s.alphavilleSeries {
		ids = append(ids, idEntry{sub.UUID})
	}
	return ids, s.found
}

func (s *dummyService) checkConnectivity() error {
	return nil
}
