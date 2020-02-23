package main

import (
	"github.com/marcsanmi/tcp-service/server/internal"
)

func main() {
	server := internal.NewServer("localhost:8001") //TODO: runtime configuration...
	server.InitListeners()
	server.Listen()
}
