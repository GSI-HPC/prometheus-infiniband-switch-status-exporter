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
