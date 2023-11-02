// -*- coding: utf-8 -*-
//
// Copyright 2023 GSI Helmholtz Centre for Heavy Ion Research
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

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
