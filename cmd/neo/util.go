package main

import (
	"os"
	"os/exec"
)

func removeFromSlice(ind int, slice []string) []string {
	return append(slice[:ind], slice[ind+1:]...)
}

func outputCmd(command string, args []string) {
	cmd := exec.Command(command, args...)
	_, err := cmd.Output()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(0)
	}
}

func runCmd(command string, args []string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		logger.Errorf("Error while running rerun command! %s", err.Error())
		os.Exit(1)
	}
}
