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
	"strings"
)

func NewSwinfoStatus() *SwinfoStatus {
	s := SwinfoStatus{}
	s.Psus = make(map[int]*SwinfoPsu)
	return &s
}

func newSwinfoPsu() *SwinfoPsu {
	return &SwinfoPsu{}
}

func NewSwinfoPsuWtihArgs(status bool, dc bool, fan bool) *SwinfoPsu {
	return &SwinfoPsu{status, dc, fan}
}

type SwinfoStatus struct {
	Psus map[int]*SwinfoPsu
	Fans bool
}

type SwinfoPsu struct {
	Status bool
	Dc     bool
	Fan    bool
}

var (
	ibswinfoPsuRegex = regexp.MustCompile(`(?m:^psu(?P<psu>[\d])\.(?P<cmp>status|dc|fan)\s*\:\s*(?P<val>OK|ERROR)$)`)
	ibswinfoPsuIndex = ibswinfoPsuRegex.SubexpIndex("psu")
	ibswinfoCmpIndex = ibswinfoPsuRegex.SubexpIndex("cmp")
	ibswinfoValIndex = ibswinfoPsuRegex.SubexpIndex("val")

	ibswinfoFansRegex = regexp.MustCompile(`(?m:^fans\s*\: (?P<val>OK|ERROR)$)`)
	ibswinfoFansIndex = ibswinfoFansRegex.SubexpIndex("val")
)

func ExtractIbswinfoStatus(input string) (*SwinfoStatus, error) {
	swStatus := NewSwinfoStatus()
	for _, line := range strings.Split(strings.TrimSuffix(input, "\n"), "\n") {

		if ibswinfoPsuRegex.MatchString(line) {
			ibswinfoPsu := ibswinfoPsuRegex.FindStringSubmatch(line)

			index, _ := strconv.Atoi(ibswinfoPsu[ibswinfoPsuIndex])
			swPsu, ok := swStatus.Psus[index]
			if ok != true {
				swPsu = newSwinfoPsu()
				swStatus.Psus[index] = swPsu
			}

			valOk := true
			if ibswinfoPsu[ibswinfoValIndex] != "OK" {
				valOk = false
			}

			comp := ibswinfoPsu[ibswinfoCmpIndex]
			if comp == "status" {
				swPsu.Status = valOk
			} else if comp == "dc" {
				swPsu.Dc = valOk
			} else if comp == "fan" {
				swPsu.Fan = valOk
			}
		} else if ibswinfoFansRegex.MatchString(line) {
			ibswinfoFans := ibswinfoFansRegex.FindStringSubmatch(line)
			if ibswinfoFans[ibswinfoFansIndex] == "OK" {
				swStatus.Fans = true
			}
		} else {
			return nil, fmt.Errorf("No regex match for line: %s", line)
		}
	}
	return swStatus, nil
}

func QueryIbswinfoStatus(lid int) (string, error) {
	lidStr := "lid-" + strconv.Itoa(lid)
	output, err := util.ExecuteCommandWithSudo("ibswinfo.sh", "-d", lidStr, "-o", "status")
	if err != nil {
		return "", fmt.Errorf("ibswinfo failed for lid %d\n"+err.Error(), lid)
	}

	return *output, nil
}
