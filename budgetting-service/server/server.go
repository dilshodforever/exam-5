package server

import (
	"log"
	"net"

	pb "github.com/dilshodforever/5-oyimtixon/genprotos/accaunts"
	budgetting "github.com/dilshodforever/5-oyimtixon/genprotos/budgets"
	category "github.com/dilshodforever/5-oyimtixon/genprotos/categories"
	goalsPb "github.com/dilshodforever/5-oyimtixon/genprotos/goals"
	transactionsPb "github.com/dilshodforever/5-oyimtixon/genprotos/transactions"
	kafkaconsumer "github.com/dilshodforever/5-oyimtixon/kafkaconsumer"
	"github.com/dilshodforever/5-oyimtixon/service"
	postgres "github.com/dilshodforever/5-oyimtixon/storage/mongo"
	"google.golang.org/grpc"
)

func Connection() net.Listener {
	// Initialize MongoDB connection
	db, err := postgres.NewMongoConnection()
	if err != nil {
		log.Fatal("Error while connecting to DB: ", err.Error())
	}

	// Set up TCP listener
	lis, err := net.Listen("tcp", ":8087")
	if err != nil {
		log.Fatal("Error while starting TCP listener: ", err.Error())
	}

	// Create a new gRPC server
	s := grpc.NewServer()

	// Initialize services
	accountService := service.NewAccountService(db)
	goalService := service.NewGoalService(db)
	transactionService := service.NewTransactionService(db)
	budget := service.NewBudgetService(db)
	categories := service.NewCategoryService(db)
	// Register services with the gRPC server
	pb.RegisterAccountServiceServer(s, accountService)
	goalsPb.RegisterGoalServiceServer(s, goalService)
	transactionsPb.RegisterTransactionServiceServer(s, transactionService)
	budgetting.RegisterBudgetServiceServer(s, budget)
	category.RegisterCategoryServiceServer(s, categories)
	// Start serving

	brokers := []string{"kafka:9092"}

	kcm := kafkaconsumer.NewKafkaConsumerManager()

	if err := kcm.RegisterConsumer(brokers, "update", "root", kafkaconsumer.UpdateBudget(budget)); err != nil {
		if err == kafkaconsumer.ErrConsumerAlreadyExists {
			log.Printf("Consumer for topic 'create-job_application' already exists")
		} else {
			log.Fatalf("Error registering consumer: %v", err)
		}
	}
	
	if err := kcm.RegisterConsumer(brokers, "CreateTransaction", "root", kafkaconsumer.CreateAccaunt(transactionService)); err != nil {
		if err == kafkaconsumer.ErrConsumerAlreadyExists {
			log.Printf("Consumer for topic 'create-job_application' already exists")
		} else {
			log.Fatalf("Error registering consumer: %v", err)
		}
	}


	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
	return lis
}
