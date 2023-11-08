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
	"fmt"
	"net/http"
	"os"
	"prometheus-infiniband-exporter/collector"
	logging "prometheus-infiniband-exporter/logging"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

const (
	defaultPort = "9860"
	version     = "0.0.2"
)

var (
	collectorCreator map[string]collector.NewCollectorHandle
	configFile       *string
)

func main() {
	printVersion := flag.Bool("version", false, "Print version")
	port := flag.String("port", defaultPort, "The port to listen on for HTTP requests")
	logLevel := flag.String("log", logging.DefaultLogLevel, "Sets log level - ERROR, WARNING, INFO, DEBUG or TRACE")
	collectSwitches := flag.Bool("collect-switches-status", false, "Enables collecting of switches status information")
	configFile = flag.String("configFile", "", "Config file to be used")

	flag.Parse()

	logging.InitLogging(*logLevel)

	if *printVersion {
		fmt.Println("Version:", version)
		os.Exit(0)
	}

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
		prometheus.MustRegister(collector(*configFile))
	}

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":"+*port, nil)
}
