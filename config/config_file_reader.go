// -*- coding: utf-8 -*-
//
// © Copyright 2023 GSI Helmholtzzentrum für Schwerionenforschung
//
// This software is distributed under
// the terms of the GNU General Public Licence version 3 (GPL Version 3),
// copied verbatim in the file "LICENCE".

package config

type ConfigFileReader struct {
	ExcludeSwitchesLids []int

	GenericConfigFileReader
}

func NewConfigFileReader(filepath string) *ConfigFileReader {
	c := new(ConfigFileReader)
	c.MustLoadFile(filepath)
	c.ExcludeSwitchesLids = c.IntList("exclude_switch_lids")
	return c
}
