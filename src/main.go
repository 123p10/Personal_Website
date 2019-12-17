package main

import (
	"github.com/Personal_Website/src/config"
	"github.com/Personal_Website/src/server"
)

func main() {
	config.LoadEnv()
	server.LoadServer()
}
