package postgres

import (
	"context"
	"fmt"
	"log"

	pb "github.com/dilshodforever/5-oyimtixon/genprotos/accaunts"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccountService struct {
	db *mongo.Database
}

func NewAccountService(db *mongo.Database) *AccountService {
	return &AccountService{db: db}
}

func (s *AccountService) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	coll := s.db.Collection("accounts")
	id := uuid.NewString()
	_, err := coll.InsertOne(ctx, bson.M{
		"id":       id,
		"UserId":   req.UserId,
		"name":     req.Name,
		"type":     req.Type,
		"balance":  0,
		"currency": req.Currency,
	})
	if err != nil {
		log.Printf("Failed to create account: %v", err)
		return &pb.CreateAccountResponse{
			Success: false,
			Message: "Failed to create account",
		}, err
	}
	return &pb.CreateAccountResponse{
		Success: true,
		Message: "Account created successfully",
	}, nil
}

func (s *AccountService) GetAccountByid(ctx context.Context, req *pb.GetByIdAccauntRequest) (*pb.GetAccountByidResponse, error) {
	coll := s.db.Collection("accounts")
	var account pb.GetAccountByidResponse
	err := coll.FindOne(ctx, bson.M{"id": req.Id}).Decode(&account)
	if err != nil {
		log.Printf("Failed to get account by id: %v", err)
		return nil, err
	}
	return &account, nil
}

func (s *AccountService) UpdateAccount(ctx context.Context, req *pb.UpdateAccountRequest) (*pb.UpdateAccountResponse, error) {
	coll := s.db.Collection("accounts")
	update := bson.M{
		"$set": bson.M{
			"name":     req.Name,
			"type":     req.Type,
			"currency": req.Currency,
		},
	}
	_, err := coll.UpdateOne(ctx, bson.M{"id": req.Id}, update)
	if err != nil {
		log.Printf("Failed to update account: %v", err)
		return &pb.UpdateAccountResponse{
			Success: false,
			Message: "Failed to update account",
		}, err
	}
	return &pb.UpdateAccountResponse{
		Success: true,
		Message: "Account updated successfully",
	}, nil
}

func (s *AccountService) DeleteAccount(ctx context.Context, req *pb.DeleteAccountRequest) (*pb.UpdateAccountResponse, error) {
	coll := s.db.Collection("accounts")
	_, err := coll.DeleteOne(ctx, bson.M{"id": req.Id})
	if err != nil {
		log.Printf("Failed to delete account: %v", err)
		return &pb.UpdateAccountResponse{
			Success: false,
			Message: "Failed to delete account",
		}, err
	}
	return &pb.UpdateAccountResponse{
		Success: true,
		Message: "Account deleted successfully",
	}, nil
}

func (s *AccountService) ListAccounts(ctx context.Context, req *pb.ListAccountsRequest) (*pb.ListAccountsResponse, error) {
	coll := s.db.Collection("accounts")
	filter := bson.M{}
	if req.Name != "" {
		filter["name"] = req.Name
	}
	if req.Type != "" {
		filter["type"] = req.Type
	}
	if req.Balance != 0 {
		filter["balance"] = req.Balance
	}
	if req.Currency != "" {
		filter["currency"] = req.Currency
	}

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		log.Printf("Failed to list accounts: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var accounts []*pb.GetAccountByidResponse
	for cursor.Next(ctx) {
		var account pb.GetAccountByidResponse
		if err := cursor.Decode(&account); err != nil {
			log.Printf("Failed to decode account: %v", err)
			return nil, err
		}
		accounts = append(accounts, &account)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}

	return &pb.ListAccountsResponse{Accounts: accounts}, nil
}

func (s *AccountService) UpdateBalance(ctx context.Context, accountID string, amount float32) error {
	coll := s.db.Collection("accounts")

	// Use the $inc operator to add the amount to the existing balance
	update := bson.M{
		"$inc": bson.M{
			"balance": amount,
		},
	}

	_, err := coll.UpdateOne(ctx, bson.M{"id": accountID}, update)
	if err != nil {
		log.Printf("Failed to update account balance: %v", err)
		return err
	}
	return nil
}

func (s *AccountService) UpdateBalanceMinus(ctx context.Context, accountID string, amount float32) error {
	coll := s.db.Collection("accounts")

	// Use the $inc operator to decrement the balance by the given amount
	update := bson.M{
		"$inc": bson.M{
			"balance": -amount,
		},
	}

	// Perform the update operation
	result, err := coll.UpdateOne(ctx, bson.M{"id": accountID}, update)
	if err != nil {
		log.Printf("Failed to update account balance: %v", err)
		return err
	}

	// Check if any document was matched by the query
	if result.MatchedCount == 0 {
		err = fmt.Errorf("no account found with ID %s", accountID)
		log.Printf("Failed to update account balance: %v", err)
		return err
	}

	return nil
}
