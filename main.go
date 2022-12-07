package main

import (
	"log"
	"os"

	"github.com/90poe/protoc-gen-gofullmethods/internal/plugin/gofullmethods"
)

func main() {
	plugin := gofullmethods.NewPlugin(os.Stdin, os.Stdout)
	if err := plugin.Run(); err != nil {
		log.Fatalf("plugin returned error, %+v", err)
	}
}
