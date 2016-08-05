package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/Financial-Times/go-fthealth/v1a"
	"github.com/Financial-Times/http-handlers-go/httphandlers"
	status "github.com/Financial-Times/service-status-go/httphandlers"
	"github.com/Financial-Times/tme-reader/tmereader"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/jawher/mow.cli"
	"github.com/rcrowley/go-metrics"
	"github.com/sethgrid/pester"
)

func init() {
	log.SetFormatter(new(log.JSONFormatter))
}

func main() {
	app := cli.App("wsod-transformer", "A RESTful API for transforming TME WSOD to UP json")
	username := app.String(cli.StringOpt{
		Name:   "tme-username",
		Value:  "",
		Desc:   "TME username used for http basic authentication",
		EnvVar: "TME_USERNAME",
	})
	password := app.String(cli.StringOpt{
		Name:   "tme-password",
		Value:  "",
		Desc:   "TME password used for http basic authentication",
		EnvVar: "TME_PASSWORD",
	})
	token := app.String(cli.StringOpt{
		Name:   "token",
		Value:  "",
		Desc:   "Token to be used for accessing TME",
		EnvVar: "TOKEN",
	})
	baseURL := app.String(cli.StringOpt{
		Name:   "base-url",
		Value:  "http://localhost:8080/transformers/wsod/",
		Desc:   "Base url",
		EnvVar: "BASE_URL",
	})
	tmeBaseURL := app.String(cli.StringOpt{
		Name:   "tme-base-url",
		Value:  "https://tme.ft.com",
		Desc:   "TME base url",
		EnvVar: "TME_BASE_URL",
	})
	port := app.Int(cli.IntOpt{
		Name:   "port",
		Value:  8080,
		Desc:   "Port to listen on",
		EnvVar: "PORT",
	})
	maxRecords := app.Int(cli.IntOpt{
		Name:   "maxRecords",
		Value:  int(10000),
		Desc:   "Maximum records to be queried to TME",
		EnvVar: "MAX_RECORDS",
	})
	slices := app.Int(cli.IntOpt{
		Name:   "slices",
		Value:  int(10),
		Desc:   "Number of requests to be executed in parallel to TME",
		EnvVar: "SLICES",
	})

	tmeTaxonomyName := app.String(cli.StringOpt{
		Name:   "tme-taxonomy-name",
		Value:  "ON",
		Desc:   "TME taxonomy name for WSOD",
		EnvVar: "TME_TAXONOMY_NAME",
	})

	app.Action = func() {
		client := getResilientClient()

		mf := new(wsodTransformer)
		s, err := newWSODService(tmereader.NewTmeRepository(client, *tmeBaseURL, *username, *password, *token, *maxRecords, *slices, *tmeTaxonomyName, &tmereader.AuthorityFiles{}, mf), *baseURL, *tmeTaxonomyName, *maxRecords)
		if err != nil {
			log.Errorf("Error while creating WSODService: [%v]", err.Error())
		}

		h := newWSODHandler(s)
		m := mux.NewRouter()

		// The top one of these feels more correct, but the lower one matches what we have in Dropwizard,
		// so it's what apps expect currently same as ping
		m.HandleFunc(status.PingPath, status.PingHandler)
		m.HandleFunc(status.PingPathDW, status.PingHandler)
		m.HandleFunc(status.BuildInfoPath, status.BuildInfoHandler)
		m.HandleFunc(status.BuildInfoPathDW, status.BuildInfoHandler)
		m.HandleFunc("/__health", v1a.Handler("WSOD Transformer Healthchecks", "Checks for accessing TME", h.HealthCheck()))
		m.HandleFunc(status.GTGPath, h.GoodToGo)

		m.HandleFunc("/transformers/wsod", h.getWSOD).Methods("GET")
		m.HandleFunc("/transformers/wsod/{uuid:([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})}", h.getWSODByUUID).Methods("GET")
		m.HandleFunc("/transformers/wsod/__ids", h.getWSODIds).Methods("GET")
		m.HandleFunc("/transformers/wsod/__count", h.getWSODCount).Methods("GET")

		http.Handle("/", m)

		log.Printf("listening on %d", *port)
		http.ListenAndServe(fmt.Sprintf(":%d", *port),
			httphandlers.HTTPMetricsHandler(metrics.DefaultRegistry,
				httphandlers.TransactionAwareRequestLoggingHandler(log.StandardLogger(), m)))
	}
	app.Run(os.Args)
}

func getResilientClient() *pester.Client {
	tr := &http.Transport{
		MaxIdleConnsPerHost: 32,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
	}
	c := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(30 * time.Second),
	}
	client := pester.NewExtendedClient(c)
	client.Backoff = pester.ExponentialBackoff
	client.MaxRetries = 5
	client.Concurrency = 1

	return client
}
