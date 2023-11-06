// -*- coding: utf-8 -*-
//
// © Copyright 2023 GSI Helmholtzzentrum für Schwerionenforschung
//
// This software is distributed under
// the terms of the GNU General Public Licence version 3 (GPL Version 3),
// copied verbatim in the file "LICENCE".

package main

import (
	"flag"
	"net/http"
	"prometheus-infiniband-exporter/collector"
	logging "prometheus-infiniband-exporter/logging"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

const (
	defaultPort = "9860"
	version     = "0.0.1"
)

var (
	collectorCreator map[string]collector.NewCollectorHandle
)

func main() {
	port := flag.String("port", defaultPort, "The port to listen on for HTTP requests")
	logLevel := flag.String("log", logging.DefaultLogLevel, "Sets log level - ERROR, WARNING, INFO, DEBUG or TRACE")
	collectSwitches := flag.Bool("collect-switches-status", false, "Enables collecting of switches status information")

	flag.Parse()

	logging.InitLogging(*logLevel)

	collectorCreator = make(map[string]collector.NewCollectorHandle)

	if *collectSwitches {
		collectorCreator[collector.SwitchesStatusCollectorName] =
			collector.NewSwitchesStatusCollector
	}

	if len(collectorCreator) == 0 {
		log.Fatalln("No collector enabled")
	}

	for name, collector := range collectorCreator {
		log.Debug("Enable collector ", name)
		prometheus.MustRegister(collector())
	}

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":"+*port, nil)
}
