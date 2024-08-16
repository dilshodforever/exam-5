package main

import (
	"github.com/dilshodforever/5-oyimtixon/kafkaconnect"
	"github.com/dilshodforever/5-oyimtixon/server"
	
)

func main() {
	server.Connection()
	kafkaconnect.ConnectToKafka()
}
