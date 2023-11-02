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

	for lid, name := range switchesIds {
		fmt.Println(lid, name)
	}

	ch <- createScrapeOkMetric(SwitchesStatusCollectorName, 1)
}

func (c *SwitchesStatusCollector) Describe(ch chan<- *prometheus.Desc) {
}
