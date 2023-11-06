// -*- coding: utf-8 -*-
//
// © Copyright 2023 GSI Helmholtzzentrum für Schwerionenforschung
//
// This software is distributed under
// the terms of the GNU General Public Licence version 3 (GPL Version 3),
// copied verbatim in the file "LICENCE".

package collector

import (
	"fmt"
	"prometheus-infiniband-exporter/ib"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

const SwitchesStatusCollectorName = "switches_status"

type SwitchesStatusCollector struct {
}

func NewSwitchesStatusCollector() prometheus.Collector {
	return &SwitchesStatusCollector{}
}

func (c *SwitchesStatusCollector) Collect(ch chan<- prometheus.Metric) {
	log.Debugln("Collecting status information of switches...")

	switchesIds, err := ib.QueryIbswitchesIds()

	if err != nil {
		log.Errorln(err)
		ch <- createScrapeOkMetric(SwitchesStatusCollectorName, 0)
		return
	}

	// for lid, name := range switchesIds {
	for lid := range switchesIds {
		_, err := ib.QueryIbswinfoStatus(lid)
		if err != nil {
			fmt.Println(err)
		}
	}

	ch <- createScrapeOkMetric(SwitchesStatusCollectorName, 1)
}

func (c *SwitchesStatusCollector) Describe(ch chan<- *prometheus.Desc) {
}
