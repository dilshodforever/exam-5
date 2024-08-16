package kafkasender

import (
	"encoding/json"
	"log"
	pbb "github.com/dilshodforever/5-oyimtixon/genprotos/budgets"
	pb "github.com/dilshodforever/5-oyimtixon/genprotos/transactions"
	"github.com/dilshodforever/5-oyimtixon/kafka"
)

func CreateTransaction(kaf kafka.KafkaProducer, request *pb.CreateTransactionRequest) (pb.TransactionResponse,error){
	response, err := json.Marshal(request)
	if err != nil {
		log.Println("cannot produce messages via kafka", err.Error())
		return pb.TransactionResponse{Message: "error"},err
	}
	err=kaf.ProduceMessages("CreateTransaction", response)
	if err != nil {
		log.Fatal("Error while ProduceMessages: ", err.Error())
		return pb.TransactionResponse{Message: "error: "},err
	}
	return pb.TransactionResponse{Message: "Success", Success: true},nil
}


func UpdateBudget(kaf kafka.KafkaProducer, request *pbb.UpdateBudgetRequest) (pbb.BudgetResponse,error){
	response, err := json.Marshal(request)
	if err != nil {
		log.Println("cannot produce messages via kafka", err.Error())
		return pbb.BudgetResponse{Message: "error"},err
	}
	err=kaf.ProduceMessages("update", response)
	if err != nil {
		log.Fatal("Error while ProduceMessages: ", err.Error())
		return pbb.BudgetResponse{Message: "error: "},err
	}
	return pbb.BudgetResponse{Message: "Success"},nil
}
