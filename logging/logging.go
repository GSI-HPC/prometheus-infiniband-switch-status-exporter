// -*- coding: utf-8 -*-
//
// © Copyright 2023 GSI Helmholtzzentrum für Schwerionenforschung
//
// This software is distributed under
// the terms of the GNU General Public Licence version 3 (GPL Version 3),
// copied verbatim in the file "LICENCE".

package logging

import (
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	DefaultLogLevel = "ERROR"
)

func InitLogging(logLevel string) {

	if logLevel == "ERROR" {
		log.SetLevel(log.ErrorLevel)
	} else if logLevel == "WARNING" {
		log.SetLevel(log.WarnLevel)
	} else if logLevel == "INFO" {
		log.SetLevel(log.InfoLevel)
	} else if logLevel == "DEBUG" {
		log.SetLevel(log.DebugLevel)
	} else if logLevel == "TRACE" {
		log.SetLevel(log.TraceLevel)
	} else {
		log.Fatal("Not supported log level set")
	}

	log.SetOutput(os.Stdout)
}
