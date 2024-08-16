package handler

import (
	pba "github.com/dilshodforever/5-oyimtixon/genprotos/accaunts"
	pbb "github.com/dilshodforever/5-oyimtixon/genprotos/budgets"
	pbc "github.com/dilshodforever/5-oyimtixon/genprotos/categories"
	pbg "github.com/dilshodforever/5-oyimtixon/genprotos/goals"
	pbn "github.com/dilshodforever/5-oyimtixon/genprotos/notifications"
	pbt "github.com/dilshodforever/5-oyimtixon/genprotos/transactions"
	"github.com/dilshodforever/5-oyimtixon/kafka"
)

type Handler struct {
	Account      pba.AccountServiceClient
	Budget       pbb.BudgetServiceClient
	Category     pbc.CategoryServiceClient
	Goal         pbg.GoalServiceClient
	Transaction  pbt.TransactionServiceClient
	Notification pbn.NotificationtServiceClient
	Redis        InMemoryStorageI
	Kafka        kafka.KafkaProducer
}

func NewHandler(
	account pba.AccountServiceClient,
	budget pbb.BudgetServiceClient,
	category pbc.CategoryServiceClient,
	goal pbg.GoalServiceClient,
	transaction pbt.TransactionServiceClient,
	notification pbn.NotificationtServiceClient,
	redis InMemoryStorageI,
    kafka kafka.KafkaProducer,

) *Handler {
	return &Handler{
		Account:      account,
		Budget:       budget,
		Category:     category,
		Goal:         goal,
		Transaction:  transaction,
		Notification: notification,
		Redis:        redis,
        Kafka: kafka,
	}
}
