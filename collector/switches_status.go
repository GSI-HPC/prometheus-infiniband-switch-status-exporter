// -*- coding: utf-8 -*-
//
// © Copyright 2023 GSI Helmholtzzentrum für Schwerionenforschung
//
// This software is distributed under
// the terms of the GNU General Public Licence version 3 (GPL Version 3),
// copied verbatim in the file "LICENCE".

package collector

import (
	"prometheus-infiniband-exporter/ib"
	"strconv"

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

	scrapeOk := 1.0

	for lid, name := range switchesIds {

		status, err := ib.QueryIbswinfoStatus(lid)

		if err != nil {
			log.Errorln(err)
			scrapeOk = 0
			continue
		}

		swinfoStatus, err := ib.ExtractIbswinfoStatus(status)

		if err != nil {
			log.Errorln(err)
			scrapeOk = 0
			continue
		}

		lidStr := strconv.Itoa(lid)

		for index := range swinfoStatus.Psus {
			psu := swinfoStatus.Psus[index]
			indexStr := strconv.Itoa(index)

			ch <- c.createPsuMetric(lidStr, name, indexStr, psu.Status,
				"power_supply_status", "Power supply status (0=ERROR, 1=OK)")

			ch <- c.createPsuMetric(lidStr, name, indexStr, psu.Dc,
				"power_supply_dc", "Power supply dc (0=ERROR, 1=OK)")

			ch <- c.createPsuMetric(lidStr, name, indexStr, psu.Fan,
				"power_supply_fan", "Power supply fan (0=ERROR, 1=OK)")
		}

		ch <- c.createFansMetric(lidStr, name, swinfoStatus.Fans,
			"fans", "Fans (0=ERROR, 1=OK)")
	}

	ch <- createScrapeOkMetric(SwitchesStatusCollectorName, scrapeOk)
}

func (c *SwitchesStatusCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (c *SwitchesStatusCollector) createPsuMetric(
	lid string, name string, index string, value bool,
	descName string, descHelp string) prometheus.Metric {

	return prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(
				Namespace,
				SwitchesStatusCollectorName,
				descName),
			descHelp,
			[]string{"lid", "name", "psu"},
			nil,
		),
		prometheus.GaugeValue,
		convertBoolToFloat(value),
		lid, name, index) // labelValues
}

func (c *SwitchesStatusCollector) createFansMetric(
	lid string, name string, value bool,
	descName string, descHelp string) prometheus.Metric {

	return prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(
				Namespace,
				SwitchesStatusCollectorName,
				descName),
			descHelp,
			[]string{"lid", "name"},
			nil,
		),
		prometheus.GaugeValue,
		convertBoolToFloat(value),
		lid, name) // labelValues
}
