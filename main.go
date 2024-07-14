package main

import (
	"blockchain-app/command"
	"os"
)

func main() {
	defer os.Exit(0)

	cmd := command.CommandLine{}
	cmd.Run()
}
