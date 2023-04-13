package main

import (
	"solace-events-consumer/src/events"
	"solace-events-consumer/src/modules/http/server"
)

func main() {
	events.ConnectReceiver(2)
	defer events.CloseReceiver()
	server.Run()
}
