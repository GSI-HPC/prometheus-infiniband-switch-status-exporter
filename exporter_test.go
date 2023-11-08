// -*- coding: utf-8 -*-
//
// © Copyright 2023 GSI Helmholtzzentrum für Schwerionenforschung
//
// This software is distributed under
// the terms of the GNU General Public Licence version 3 (GPL Version 3),
// copied verbatim in the file "LICENCE".

package main

import (
	"prometheus-infiniband-exporter/config"
	"prometheus-infiniband-exporter/ib"
	"prometheus-infiniband-exporter/util"
	"testing"
)

var (
	testConfigFile     = "config_example.yml"
	testIbswitchesFile = "ibswitches.txt"
	testIbswinfoFile   = "ibswinfo.txt"
)

func TestQueryIbswitches(t *testing.T) {
	var data ib.IbswitchesIds
	var err error

	_, err = ib.ExtractIbswitchesIds("no valid")

	if err == nil {
		t.Errorf("ExtractSwitchesId must return an error on no valid data")
	}

	output := util.MustReadFile(&testIbswitchesFile)

	data, err = ib.ExtractIbswitchesIds(output)

	if err != nil {
		t.Errorf("ExtractSwitchesId no error expected, got: %s", err)
	}

	expectedLen := 8
	receivedLen := len(data)
	if expectedLen != receivedLen {
		t.Errorf("Incomplete count of switches - expected: %d - received: %d",
			expectedLen, receivedLen)
	}

	expectedLid := 770
	expectedName := "Quantum Mellanox Technologies"
	receivedName, keyFound := data[expectedLid]

	if keyFound == false {
		t.Error("lid not found for switch", expectedLid)
	}

	if expectedName != receivedName {
		t.Errorf("Switch name not found for lid %d, expected: %s, received: %s",
			expectedLid, expectedName, receivedName)
	}
}

func TestIbswinfoSwitchStatus(t *testing.T) {

	expectedSwinfoStatus := ib.NewSwinfoStatus()
	expectedSwinfoStatus.Psus[0] = ib.NewSwinfoPsuWtihArgs(true, true, true)
	expectedSwinfoStatus.Psus[1] = ib.NewSwinfoPsuWtihArgs(false, false, false)
	expectedSwinfoStatus.Fans = true

	output := util.MustReadFile(&testIbswinfoFile)
	receivedSwinfoStatus, _ := ib.ExtractIbswinfoStatus(output)

	for psuIndex, expectedPsuInfo := range expectedSwinfoStatus.Psus {
		if expectedPsuInfo.Status != receivedSwinfoStatus.Psus[psuIndex].Status {
			t.Errorf("PSU[%d] status failed, expected: %t, received: %t",
				psuIndex,
				expectedPsuInfo.Status,
				receivedSwinfoStatus.Psus[psuIndex].Status)
		}
		if expectedPsuInfo.Dc != receivedSwinfoStatus.Psus[psuIndex].Dc {
			t.Errorf("PSU[%d] dc failed, expected: %t, received: %t",
				psuIndex,
				expectedPsuInfo.Dc,
				receivedSwinfoStatus.Psus[psuIndex].Dc)
		}
		if expectedPsuInfo.Fan != receivedSwinfoStatus.Psus[psuIndex].Fan {
			t.Errorf("PSU[%d] fan failed, expected: %t, received: %t",
				psuIndex,
				expectedPsuInfo.Fan,
				receivedSwinfoStatus.Psus[psuIndex].Fan)
		}
	}

	if expectedSwinfoStatus.Fans != receivedSwinfoStatus.Fans {
		t.Errorf("Test for fans failed, expected: %t, received: %t", expectedSwinfoStatus.Fans, receivedSwinfoStatus.Fans)
	}
}

func TestSwitchesLids(t *testing.T) {
	configFileReader := config.NewConfigFileReader(testConfigFile)

	expectedLenSwitchesLids := 8
	receivedLenSwitchesLids := len(configFileReader.ExcludeSwitchesLids)
	if expectedLenSwitchesLids != receivedLenSwitchesLids {
		t.Errorf("Count of switches lids failed, expected: %d, received: %d",
			expectedLenSwitchesLids, receivedLenSwitchesLids)
	}
}
