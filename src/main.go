package main

import (
	"os"

	"github.com/Personal_Website/src/config"
	"github.com/Personal_Website/src/server"
)

func main() {
	if os.Getenv("PRODUCTION") != "TRUE" {
		config.LoadEnv()
	}
	server.LoadServer()
}
