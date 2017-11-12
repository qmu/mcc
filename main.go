package main

import (
	"log"
	"os"

	"github.com/qmu/mcc/controller"
)

func main() {
	if err := controller.Run(); err != nil {
		log.Panicln(err)
		os.Exit(1)
	}
}
