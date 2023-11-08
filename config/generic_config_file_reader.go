// -*- coding: utf-8 -*-
//
// © Copyright 2023 GSI Helmholtzzentrum für Schwerionenforschung
//
// This software is distributed under
// the terms of the GNU General Public Licence version 3 (GPL Version 3),
// copied verbatim in the file "LICENCE".

package config

import (
	"log"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yamlv3"
)

type GenericConfigFileReader struct{}

func (*GenericConfigFileReader) MustLoadFile(filepath string) {
	config.AddDriver(yamlv3.Driver)

	err := config.LoadFiles(filepath)
	if err != nil {
		panic(err)
	}
}

func (*GenericConfigFileReader) MustHaveString(key string) string {
	value := config.String(key)

	if len(value) == 0 {
		log.Panic("Key not found or has no value in config file: ", key)
	}

	return value
}

func (*GenericConfigFileReader) MustHaveStringList(key string) []string {
	list := config.Strings(key)

	if len(list) == 0 {
		log.Panic("Key not found or has no list items in config file: ", key)
	}

	return list
}

func (*GenericConfigFileReader) MustHaveMap(key string) map[string]string {
	mmap := config.StringMap(key)

	if len(mmap) == 0 {
		log.Panic("Map not found or is empty: ", key)
	}

	return mmap
}

func (*GenericConfigFileReader) IntList(key string) []int {
	return config.Ints(key)
}
