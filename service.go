package main

import (
	"net/http"

	"github.com/Financial-Times/tme-reader/tmereader"
	log "github.com/Sirupsen/logrus"
)

type httpClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

type alphavilleSeriesService interface {
	getAlphavilleSeries() ([]alphavilleSeriesLink, bool)
	getAlphavilleSeriesIds() ([]idEntry, bool)
	getAlphavilleSeriesByUUID(uuid string) (alphavilleSeries, bool)
	getAlphavilleSeriesCount() int
	checkConnectivity() error
}

type alphavilleSeriesServiceImpl struct {
	repository            tmereader.Repository
	baseURL               string
	alphavilleSeriesMap   map[string]alphavilleSeries
	alphavilleSeriesLinks []alphavilleSeriesLink
	taxonomyName          string
	maxTmeRecords         int
}

func newAlphavilleSeriesService(repo tmereader.Repository, baseURL string, taxonomyName string, maxTmeRecords int) (alphavilleSeriesService, error) {
	s := &alphavilleSeriesServiceImpl{repository: repo, baseURL: baseURL, taxonomyName: taxonomyName, maxTmeRecords: maxTmeRecords}
	err := s.init()
	if err != nil {
		return &alphavilleSeriesServiceImpl{}, err
	}
	return s, nil
}

func (s *alphavilleSeriesServiceImpl) init() error {
	s.alphavilleSeriesMap = make(map[string]alphavilleSeries)
	responseCount := 0
	log.Printf("Fetching Alphaville Series from TME\n")
	for {
		terms, err := s.repository.GetTmeTermsFromIndex(responseCount)
		if err != nil {
			return err
		}

		if len(terms) < 1 {
			log.Printf("Finished fetching Alphaville Series from TME\n")
			break
		}
		s.initAlphavilleSeriesMap(terms)
		responseCount += s.maxTmeRecords
	}
	log.Printf("Added %d Alphaville Series links\n", len(s.alphavilleSeriesLinks))

	return nil
}

func (s *alphavilleSeriesServiceImpl) getAlphavilleSeries() ([]alphavilleSeriesLink, bool) {
	if len(s.alphavilleSeriesLinks) > 0 {
		return s.alphavilleSeriesLinks, true
	}
	return s.alphavilleSeriesLinks, false
}

func (s *alphavilleSeriesServiceImpl) getAlphavilleSeriesIds() ([]idEntry, bool) {
	if len(s.alphavilleSeriesMap) > 0 {
		ids := make([]idEntry, len(s.alphavilleSeriesMap))
		i := 0
		for k := range s.alphavilleSeriesMap {
			ids[i] = idEntry{k}
			i++
		}
		return ids, true
	}

	return make([]idEntry, 0), false
}

func (s *alphavilleSeriesServiceImpl) getAlphavilleSeriesCount() int {
	return len(s.alphavilleSeriesMap)
}

func (s *alphavilleSeriesServiceImpl) getAlphavilleSeriesByUUID(uuid string) (alphavilleSeries, bool) {
	alphavilleSeries, found := s.alphavilleSeriesMap[uuid]
	return alphavilleSeries, found
}

func (s *alphavilleSeriesServiceImpl) checkConnectivity() error {
	// TODO: Can we just hit an endpoint to check if TME is available? Or do we need to make sure we get genre taxonmies back? Maybe a healthcheck or gtg endpoint?
	// TODO: Can we use a count from our responses while actually in use to trigger a healthcheck?
	//	_, err := s.repository.GetTmeTermsFromIndex(1)
	//	if err != nil {
	//		return err
	//	}
	return nil
}

func (s *alphavilleSeriesServiceImpl) initAlphavilleSeriesMap(terms []interface{}) {
	for _, iTerm := range terms {
		t := iTerm.(term)
		top := transformAlphavilleSeries(t, s.taxonomyName)
		s.alphavilleSeriesMap[top.UUID] = top
		s.alphavilleSeriesLinks = append(s.alphavilleSeriesLinks, alphavilleSeriesLink{APIURL: s.baseURL + top.UUID})
	}
}
