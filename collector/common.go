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

import "github.com/prometheus/client_golang/prometheus"

const Namespace = "infiniband"

// Function signature for NewCollector...
type NewCollectorHandle func() prometheus.Collector

type metricTemplate struct {
	desc         *prometheus.Desc
	valueType    prometheus.ValueType
	valueCreator func(string) (float64, error)
}

func createScrapeOkMetric(collector string, value float64) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "collector", "scrape_ok"),
			"Indicates if scrape of a collector was OK",
			[]string{"name"},
			nil,
		),
		prometheus.GaugeValue,
		value,
		collector,
	)
}
