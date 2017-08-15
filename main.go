package main

import (
	"log"
	"os"

	_ "github.com/k0kubun/pp"
	"github.com/qmu/mcc/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Panicln(err)
		os.Exit(1)
	}
}
