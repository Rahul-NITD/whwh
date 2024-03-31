package main

import (
	"os"

	"github.com/aargeee/whwh/systems/client/cli"
)

func main() {
	cli.NewCLI(os.Stdout).BeginCLI()
}
