package main

import (
	"github.com/utilyre/session-auth/application"
	"github.com/utilyre/session-auth/config"
)

func main() {
	cfg := config.Load()
	application.New(cfg).Setup().Start()
}
