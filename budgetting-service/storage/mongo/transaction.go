package postgres

import (
	"context"
	"log"

	"github.com/dilshodforever/5-oyimtixon/genprotos/transactions"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *AccountService) CreateTransaction(ctx context.Context, req *transactions.CreateTransactionRequest) (*transactions.TransactionResponse, error) {
	coll := s.db.Collection("transactions")
	id := uuid.NewString()
	_, err := coll.InsertOne(ctx, bson.M{
		"id":          id,
		"UserId":      req.UserId,
		"AccountId":   req.AccountId,
		"CategoryId":  req.CategoryId,
		"amount":      req.Amount,
		"type":        req.Type,
		"description": req.Description,
		"date":        req.Date,
	})
	//
	if err != nil {
		log.Printf("Failed to create transaction: %v", err)
		return &transactions.TransactionResponse{Success: false, Message: "Failed to create transaction"}, err
	}
	return &transactions.TransactionResponse{Success: true, Message: "Transaction passed successfully"}, nil
}

func (s *AccountService) GetTransaction(ctx context.Context, req *transactions.GetTransactionRequest) (*transactions.GetTransactionResponse, error) {
	coll := s.db.Collection("transactions")

	var transaction transactions.GetTransactionResponse
	err := coll.FindOne(ctx, bson.M{"id": req.Id}).Decode(&transaction)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("No transaction found with id: %v", req.Id)
			return nil, err
		}
		log.Printf("Failed to get transaction by id: %v", err)
		return nil, err
	}

	return &transaction, nil
}

func (s *AccountService) UpdateTransaction(ctx context.Context, req *transactions.UpdateTransactionRequest) (*transactions.TransactionResponse, error) {
	coll := s.db.Collection("transactions")

	update := bson.M{}
	if req.AccountId != "" {
		update["AccountId"] = req.AccountId
	}
	if req.CategoryId != "" {
		update["CategoryId"] = req.CategoryId
	}
	if req.Amount > 0 {
		update["amount"] = req.Amount
	}
	if req.Type != "" {
		update["type"] = req.Type
	}
	if req.Description != "" {
		update["description"] = req.Description
	}
	if req.Date != "" {
		update["date"] = req.Date
	}

	if len(update) == 0 {
		return &transactions.TransactionResponse{Success: false, Message: "Nothing to update"}, nil
	}

	_, err := coll.UpdateOne(ctx, bson.M{"id": req.Id}, bson.M{"$set": update})
	if err != nil {
		log.Printf("Failed to update transaction: %v", err)
		return &transactions.TransactionResponse{Success: false, Message: "Failed to update transaction"}, err
	}

	return &transactions.TransactionResponse{Success: true, Message: "Transaction updated successfully"}, nil
}

func (s *AccountService) DeleteTransaction(ctx context.Context, req *transactions.DeleteTransactionRequest) (*transactions.TransactionResponse, error) {
	coll := s.db.Collection("transactions")

	_, err := coll.DeleteOne(ctx, bson.M{"id": req.Id})
	if err != nil {
		log.Printf("Failed to delete transaction: %v", err)
		return &transactions.TransactionResponse{Success: false, Message: "Failed to delete transaction"}, err
	}

	return &transactions.TransactionResponse{Success: true, Message: "Transaction deleted successfully"}, nil
}

func (s *AccountService) ListTransactions(ctx context.Context, req *transactions.ListTransactionsRequest) (*transactions.ListTransactionsResponse, error) {
	coll := s.db.Collection("transactions")

	filter := bson.M{}
	if req.AccountId != "" {
		filter["AccountId"] = req.AccountId
	}
	if req.CategoryId != "" {
		filter["CategoryId"] = req.CategoryId
	}
	if req.Amount > 0 {
		filter["amount"] = req.Amount
	}
	if req.Type != "" {
		filter["type"] = req.Type
	}
	if req.Description != "" {
		filter["description"] = req.Description
	}
	if req.Date != "" {
		filter["date"] = req.Date
	}

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		log.Printf("Failed to list transactions: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactionsList []*transactions.GetTransactionResponse
	for cursor.Next(ctx) {
		var transaction transactions.GetTransactionResponse
		if err := cursor.Decode(&transaction); err != nil {
			log.Printf("Failed to decode transaction: %v", err)
			return nil, err
		}
		transactionsList = append(transactionsList, &transaction)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}

	return &transactions.ListTransactionsResponse{Transactions: transactionsList}, nil
}
