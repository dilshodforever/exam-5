package kafkasender

import (
	"encoding/json"
	"log"

	"github.com/dilshodforever/5-oyimtixon/model"
	"github.com/dilshodforever/5-oyimtixon/kafka"
)

func CreateNotification(kaf kafka.KafkaProducer, request model.Send) error{
	response, err := json.Marshal(request)
	if err != nil {
		log.Println("cannot produce messages via kafka", err.Error())
		return err
	}
	err=kaf.ProduceMessages("create", response)
	if err != nil {
		log.Fatal("Error while ProduceMessages: ", err.Error())
		return err
	}
	return nil
}
