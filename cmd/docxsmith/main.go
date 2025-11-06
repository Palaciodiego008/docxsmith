package main

import (
	"os"

	"github.com/Palaciodiego008/docxsmith/internal/cli"
)

func main() {
	// Pass all arguments except the program name to the CLI
	cli.Run(os.Args[1:])
}
