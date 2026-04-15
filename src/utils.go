package main

import (
	"os"
	"os/exec"
)

func isAdmin() bool {
	if OS == "linux" {
		return os.Getenv("SUDO_USER") != ""
	}

	return exec.Command("net", "session").Run() == nil
}

func countPorts(input []int, match int) int {
	count := 0
	for _, val := range input {
		if val == match {
			count++
		}
	}
	return count
}
