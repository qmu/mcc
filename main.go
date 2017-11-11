package main

import (
	"log"
	"os"

	"github.com/qmu/mcc/commands"
)

func main() {
	if err := commands.Mcc(); err != nil {
		log.Panicln(err)
		os.Exit(1)
	}
}
