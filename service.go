package main

import (
	"net/http"

	"github.com/Financial-Times/tme-reader/tmereader"
	log "github.com/Sirupsen/logrus"
)

type httpClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

type wsodService interface {
	getWSOD() ([]wsodLink, bool)
	getWSODIds() ([]idEntry, bool)
	getWSODByUUID(uuid string) (wsod, bool)
	getWSODCount() int
	checkConnectivity() error
}

type wsodServiceImpl struct {
	repository    tmereader.Repository
	baseURL       string
	wsodMap       map[string]wsod
	wsodLinks     []wsodLink
	taxonomyName  string
	maxTmeRecords int
}

func newWSODService(repo tmereader.Repository, baseURL string, taxonomyName string, maxTmeRecords int) (wsodService, error) {
	s := &wsodServiceImpl{repository: repo, baseURL: baseURL, taxonomyName: taxonomyName, maxTmeRecords: maxTmeRecords}
	err := s.init()
	if err != nil {
		return &wsodServiceImpl{}, err
	}
	return s, nil
}

func (s *wsodServiceImpl) init() error {
	s.wsodMap = make(map[string]wsod)
	responseCount := 0
	log.Printf("Fetching WSOD from TME\n")
	for {
		terms, err := s.repository.GetTmeTermsFromIndex(responseCount)
		if err != nil {
			return err
		}

		if len(terms) < 1 {
			log.Printf("Finished fetching WSOD from TME\n")
			break
		}
		s.initWSODMap(terms)
		responseCount += s.maxTmeRecords
	}
	log.Printf("Added %d WSOD links\n", len(s.wsodLinks))

	return nil
}

func (s *wsodServiceImpl) getWSOD() ([]wsodLink, bool) {
	if len(s.wsodLinks) > 0 {
		return s.wsodLinks, true
	}
	return s.wsodLinks, false
}

func (s *wsodServiceImpl) getWSODIds() ([]idEntry, bool) {
	if len(s.wsodMap) > 0 {
		ids := make([]idEntry, len(s.wsodMap))
		i := 0
		for k := range s.wsodMap {
			ids[i] = idEntry{k}
			i++
		}
		return ids, true
	}

	return make([]idEntry, 0), false
}

func (s *wsodServiceImpl) getWSODCount() int {
	return len(s.wsodMap)
}

func (s *wsodServiceImpl) getWSODByUUID(uuid string) (wsod, bool) {
	wsodID, found := s.wsodMap[uuid]
	return wsodID, found
}

func (s *wsodServiceImpl) checkConnectivity() error {
	// TODO: Can we just hit an endpoint to check if TME is available? Or do we need to make sure we get genre taxonmies back? Maybe a healthcheck or gtg endpoint?
	// TODO: Can we use a count from our responses while actually in use to trigger a healthcheck?
	//	_, err := s.repository.GetTmeTermsFromIndex(1)
	//	if err != nil {
	//		return err
	//	}
	return nil
}

func (s *wsodServiceImpl) initWSODMap(terms []interface{}) {
	for _, iTerm := range terms {
		t := iTerm.(term)
		top := transformWSOD(t, s.taxonomyName)
		s.wsodMap[top.UUID] = top
		s.wsodLinks = append(s.wsodLinks, wsodLink{APIURL: s.baseURL + top.UUID})
	}
}
