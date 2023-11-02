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

package util

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

// Reads a file and panics on error
func MustReadFile(file *string) string {
	data, err := os.ReadFile(*file)

	if err != nil {
		panic(err)
	}

	return string(data)
}

func ExecuteCommandWithSudo(command string, args ...string) (*string, error) {
	cmdWithArgs := append([]string{command}, args...)

	cmd := exec.Command("sudo", cmdWithArgs...)

	log.Debug("Executing command: ", cmd.String())

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		errMsg := fmt.Sprintf("Error: %s", err.Error())
		if len(stderr.String()) != 0 {
			errMsg += fmt.Sprintf(" - STDERR: %s", stderr.String())
		}
		if len(stdout.String()) != 0 {
			errMsg += fmt.Sprintf(" - STDOUT: %s", stdout.String())
		}
		return nil, fmt.Errorf(errMsg)
	}

	// TrimSpace on []byte is more efficient than TrimSpace on a string since it creates a copy
	content := string(bytes.TrimSpace(stdout.Bytes()))

	if len(content) == 0 {
		return nil, errors.New("Empty content recieved for command: " + cmd.String())
	}

	return &content, nil
}
