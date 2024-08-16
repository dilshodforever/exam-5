package main

import (
	"fmt"
	"log"

	"github.com/dilshodforever/5-oyimtixon/server"
)

func main() {
	r, err:=server.Connect()
	if err!=nil{
		panic(err)
	}
	fmt.Println("Server started on port:8081")
	err = r.Run(":8081")
	if err != nil {
		log.Fatal("Error while running server: ", err.Error())
	}
}
