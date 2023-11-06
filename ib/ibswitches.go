// -*- coding: utf-8 -*-
//
// © Copyright 2023 GSI Helmholtzzentrum für Schwerionenforschung
//
// This software is distributed under
// the terms of the GNU General Public Licence version 3 (GPL Version 3),
// copied verbatim in the file "LICENCE".

package ib

import (
	"fmt"
	"prometheus-infiniband-exporter/util"
	"regexp"
	"strconv"
)

var (
	ibswitchesRegex     = regexp.MustCompile(`(?m:^Switch.* ports [\d]+ \"(?P<name>.*)\" .* lid (?P<lid>[\d]+))`)
	ibswitchesNameIndex = ibswitchesRegex.SubexpIndex("name")
	ibswitchesLidIndex  = ibswitchesRegex.SubexpIndex("lid")
)

type IbswitchesIds map[int]string // key=lid, value=name

// Public for testing
func ExtractIbswitchesIds(input string) (IbswitchesIds, error) {
	ibswitchesOverview := make(IbswitchesIds)

	ibswitches := ibswitchesRegex.FindAllStringSubmatch(input, -1)

	if ibswitches == nil {
		return nil, fmt.Errorf("No switches found")
	}

	for _, ibswitch := range ibswitches {
		lid, err := strconv.Atoi(ibswitch[ibswitchesLidIndex])

		if err != nil {
			return nil, err
		}

		ibswitchesOverview[lid] = ibswitch[ibswitchesNameIndex]
	}

	return ibswitchesOverview, nil
}

func QueryIbswitchesIds() (IbswitchesIds, error) {

	output, err := util.ExecuteCommandWithSudo("ibswitches")
	if err != nil {
		return nil, err
	}

	ibswitchesIds, err := ExtractIbswitchesIds(*output)
	if err != nil {
		return nil, err
	}

	return ibswitchesIds, nil
}
