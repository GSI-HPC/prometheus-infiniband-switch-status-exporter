// -*- coding: utf-8 -*-
//
// © Copyright 2023 GSI Helmholtzzentrum für Schwerionenforschung
//
// This software is distributed under
// the terms of the GNU General Public Licence version 3 (GPL Version 3),
// copied verbatim in the file "LICENCE".

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
