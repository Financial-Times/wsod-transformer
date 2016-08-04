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
const getWSODResponse = `[{"apiUrl":"http://localhost:8080/transformers/wsod/bba39990-c78d-3629-ae83-808c333c6dbc"}]` + "\n"
const getWSODByUUIDResponse = `{"uuid":"bba39990-c78d-3629-ae83-808c333c6dbc","alternativeIdentifiers":{"TME":["MTE3-U3ViamVjdHM="],"uuids":["bba39990-c78d-3629-ae83-808c333c6dbc"]},"prefLabel":"Global WSOD","type":"WSOD"}` + "\n"

func TestHandlers(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name         string
		req          *http.Request
		dummyService wsodService
		statusCode   int
		contentType  string // Contents of the Content-Type header
		body         string
	}{
		{"Success - get wsod by uuid", newRequest("GET", fmt.Sprintf("/transformers/wsod/%s", testUUID)), &dummyService{found: true, wsod: []wsod{getDummyWSOD(testUUID, "Global WSOD", "MTE3-U3ViamVjdHM=")}}, http.StatusOK, "application/json", getWSODByUUIDResponse},
		{"Not found - get wsod by uuid", newRequest("GET", fmt.Sprintf("/transformers/wsod/%s", testUUID)), &dummyService{found: false, wsod: []wsod{wsod{}}}, http.StatusNotFound, "application/json", ""},
		{"Success - get wsod", newRequest("GET", "/transformers/wsod"), &dummyService{found: true, wsod: []wsod{wsod{UUID: testUUID}}}, http.StatusOK, "application/json", getWSODResponse},
		{"Not found - get wsod", newRequest("GET", "/transformers/wsod"), &dummyService{found: false, wsod: []wsod{}}, http.StatusNotFound, "application/json", ""},
		{"Success - get wsod ids", newRequest("GET", "/transformers/wsod/__ids"), &dummyService{found: true, wsod: []wsod{wsod{UUID: testUUID}}}, http.StatusOK, "", getIdsResponse},
		{"Success - get 2 wsod ids", newRequest("GET", "/transformers/wsod/__ids"), &dummyService{found: true, wsod: []wsod{wsod{UUID: oneUUID}, wsod{UUID: twoUUID}}}, http.StatusOK, "", getTwoIdsResponse},
		{"Not found - get wsod", newRequest("GET", "/transformers/wsod/__ids"), &dummyService{found: false, wsod: []wsod{}}, http.StatusNotFound, "application/json", ""},
		{"Success - get wsod count", newRequest("GET", "/transformers/wsod/__count"), &dummyService{found: true, wsod: []wsod{wsod{UUID: testUUID}}}, http.StatusOK, "application/json", countResponse},
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

func router(s wsodService) *mux.Router {
	m := mux.NewRouter()
	h := newWSODHandler(s)
	m.HandleFunc("/transformers/wsod", h.getWSOD).Methods("GET")
	m.HandleFunc("/transformers/wsod/{uuid:([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})}", h.getWSODByUUID).Methods("GET")
	m.HandleFunc("/transformers/wsod/__ids", h.getWSODIds).Methods("GET")
	m.HandleFunc("/transformers/wsod/__count", h.getWSODCount).Methods("GET")
	return m
}

type dummyService struct {
	found bool
	wsod  []wsod
}

func (s *dummyService) getWSOD() ([]wsodLink, bool) {
	var wsodLinks []wsodLink
	for _, sub := range s.wsod {
		wsodLinks = append(wsodLinks, wsodLink{APIURL: "http://localhost:8080/transformers/wsod/" + sub.UUID})
	}
	return wsodLinks, s.found
}

func (s *dummyService) getWSODByUUID(uuid string) (wsod, bool) {
	return s.wsod[0], s.found
}

func (s *dummyService) getWSODCount() int {
	return len(s.wsod)
}

func (s *dummyService) getWSODIds() ([]idEntry, bool) {
	var ids []idEntry
	for _, sub := range s.wsod {
		ids = append(ids, idEntry{sub.UUID})
	}
	return ids, s.found
}

func (s *dummyService) checkConnectivity() error {
	return nil
}
