package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAlphavilleSeries(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name             string
		baseURL          string
		terms            []term
		alphavilleSeries []alphavilleSeriesLink
		found            bool
		err              error
	}{
		{"Success", "localhost:8080/transformers/alphavilleseries/",
			[]term{term{CanonicalName: "Z_Archive", RawID: "b8337559-ac08-3404-9025-bad51ebe2fc7"}, term{CanonicalName: "Feature", RawID: "mNGQ2MWQ0NDMtMDc5Mi00NWExLTlkMGQtNWZhZjk0NGExOWU2-Z2VucVz"}},
			[]alphavilleSeriesLink{alphavilleSeriesLink{APIURL: "localhost:8080/transformers/alphavilleseries/41c03fd4-8f24-3130-9f20-4d25c0909594"},
				alphavilleSeriesLink{APIURL: "localhost:8080/transformers/alphavilleseries/44dc1ad7-76f1-39be-8ff1-3d5da91520ee"}}, true, nil},
		{"Error on init", "localhost:8080/transformers/alphavilleseries/", []term{}, []alphavilleSeriesLink(nil), false, errors.New("Error getting taxonomy")},
	}

	for _, test := range tests {
		repo := dummyRepo{terms: test.terms, err: test.err}
		service, err := newAlphavilleSeriesService(&repo, test.baseURL, "Series", 10000)
		expectedAlphavilleSeries, found := service.getAlphavilleSeries()
		assert.Equal(test.alphavilleSeries, expectedAlphavilleSeries, fmt.Sprintf("%s: Expected Alphaville Series link incorrect", test.name))
		assert.Equal(test.found, found)
		assert.Equal(test.err, err)
	}
}

func TestGetAlphavilleSeriesByUuid(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name             string
		terms            []term
		uuid             string
		alphavilleSeries alphavilleSeries
		found            bool
		err              error
	}{
		{"Success", []term{term{CanonicalName: "Z_Archive", RawID: "b8337559-ac08-3404-9025-bad51ebe2fc7"}, term{CanonicalName: "Feature", RawID: "TkdRMk1XUTBORE10TURjNU1pMDBOV0V4TFRsa01HUXROV1poWmprME5HRXhPV1UyLVoyVnVjbVZ6-U2VjdGlvbnM=]"}},
			"41c03fd4-8f24-3130-9f20-4d25c0909594", getDummyAlphavilleSeries("41c03fd4-8f24-3130-9f20-4d25c0909594", "Z_Archive", "YjgzMzc1NTktYWMwOC0zNDA0LTkwMjUtYmFkNTFlYmUyZmM3-U2VyaWVz"), true, nil},
		{"Not found", []term{term{CanonicalName: "Z_Archive", RawID: "845dc7d7-ae89-4fed-a819-9edcbb3fe507"}, term{CanonicalName: "Feature", RawID: "NGQ2MWdefsdfsfcmVz"}},
			"some uuid", alphavilleSeries{}, false, nil},
		{"Error on init", []term{}, "some uuid", alphavilleSeries{}, false, errors.New("Error getting taxonomy")},
	}

	for _, test := range tests {
		repo := dummyRepo{terms: test.terms, err: test.err}
		service, err := newAlphavilleSeriesService(&repo, "", "Series", 10000)
		expectedAlphavilleSeries, found := service.getAlphavilleSeriesByUUID(test.uuid)
		assert.Equal(test.alphavilleSeries, expectedAlphavilleSeries, fmt.Sprintf("%s: Expected Alphaville Series incorrect", test.name))
		assert.Equal(test.found, found)
		assert.Equal(test.err, err)
	}
}

type dummyRepo struct {
	terms []term
	err   error
}

func (d *dummyRepo) GetTmeTermsFromIndex(startRecord int) ([]interface{}, error) {
	if startRecord > 0 {
		return nil, d.err
	}
	var interfaces = make([]interface{}, len(d.terms))
	for i, data := range d.terms {
		interfaces[i] = data
	}
	return interfaces, d.err
}

func (d *dummyRepo) GetTmeTermById(uuid string) (interface{}, error) {
	return d.terms[0], d.err
}

func getDummyAlphavilleSeries(uuid string, prefLabel string, tmeID string) alphavilleSeries {
	return alphavilleSeries{
		UUID:      uuid,
		PrefLabel: prefLabel,
		Type:      "AlphavilleSeries",
		AlternativeIdentifiers: alternativeIdentifiers{TME: []string{tmeID}, Uuids: []string{uuid}}}
}
