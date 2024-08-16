package kafka

import (
	"context"
	"encoding/json"
	"log"

	pb "github.com/dilshodforever/5-oyimtixon/genprotos/transactions"
	pbb "github.com/dilshodforever/5-oyimtixon/genprotos/budgets"
	"github.com/dilshodforever/5-oyimtixon/service"
)

func CreateAccaunt(rootService *service.TransactionService) func(message []byte) {
	return func(message []byte) {
		var app pb.CreateTransactionRequest
		if err := json.Unmarshal(message, &app); err != nil {
			log.Printf("Cannot unmarshal JSON: %v", err)
			return
		}
		_, err := rootService.CreateTransaction(context.Background(),&app)
		if err != nil {
			log.Printf("Cannot create evaluation via Kafka: %v", err)
			return
		}
		log.Print("Created evaluation")
	}
}


func UpdateBudget(rootService *service.BudgetService) func(message []byte) {
	return func(message []byte) {
		var app pbb.CreateBudgetRequest
		if err := json.Unmarshal(message, &app); err != nil {
			log.Printf("Cannot unmarshal JSON: %v", err)
			return
		}
		_, err := rootService.CreateBudget(context.Background(),&app)
		if err != nil {
			log.Printf("Cannot create evaluation via Kafka: %v", err)
			return
		}
		log.Print("Created evaluation")
	}
}



