package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/dilshodforever/5-oyimtixon/config"
	u "github.com/dilshodforever/5-oyimtixon/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	Db    *mongo.Database
	Accs  u.AccountService
	Budg  u.BudgetService
	Cats  u.CategoryService
	Goals u.GoalService
	Trans u.TransactionService
}


func NewMongoConnection() (u.InitRoot, error) {
	config:=config.Load()
	clientOptions := options.Client().ApplyURI(config.MongoDBConnection)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Error: Couldn't connect to the database.", err)
	}

	fmt.Println("Connected to MongoDB!")

	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("Budgetting")

	return &MongoStorage{Db: db}, err
}

func (s *MongoStorage) Account() u.AccountService {
	if s.Accs == nil {
		s.Accs = &AccountService{db: s.Db} // Ensure AccountService implements u.AccountService
	}
	return s.Accs
}

func (s *MongoStorage) Budget() u.BudgetService {
	if s.Budg == nil {
		s.Budg = &AccountService{db: s.Db} // Ensure BudgetStorage implements u.BudgetService
	}
	return s.Budg
}

func (s *MongoStorage) Category() u.CategoryService {
	if s.Cats == nil {
		s.Cats = &AccountService{db: s.Db} // Ensure AccountService implements u.CategoryService
	}
	return s.Cats
}

func (s *MongoStorage) Goal() u.GoalService {
	if s.Goals == nil {
		s.Goals = &AccountService{db: s.Db} // Ensure AccountService implements u.GoalService
	}
	return s.Goals
}

func (s *MongoStorage) Transaction() u.TransactionService {
	if s.Trans == nil {
		s.Trans = &AccountService{db: s.Db} // Ensure AccountService implements u.TransactionService
	}
	return s.Trans
}

