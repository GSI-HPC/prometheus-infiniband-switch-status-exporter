// -*- coding: utf-8 -*-
//
// © Copyright 2023 GSI Helmholtzzentrum für Schwerionenforschung
//
// This software is distributed under
// the terms of the GNU General Public Licence version 3 (GPL Version 3),
// copied verbatim in the file “LICENCE”.

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
	cmd := exec.Command(command, args...)

	uid := os.Getuid()
	if uid != 0 {
		cmdWithArgs := append([]string{command}, args...)
		cmd = exec.Command("sudo", cmdWithArgs...)
	}

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
