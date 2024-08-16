package postgres

import (
	"context"
	"log"
	"time"

	"github.com/dilshodforever/5-oyimtixon/genprotos/budgets"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *AccountService) CreateBudget(ctx context.Context, req *budgets.CreateBudgetRequest) (*budgets.BudgetResponse, error) {
	coll := s.db.Collection("budgets")
	id := uuid.NewString()
	_, err := coll.InsertOne(ctx, bson.M{
		"id":         id,
		"UserId":     req.UserId,
		"CategoryId": req.CategoryId,
		"amount":     req.Amount,
		"period":     req.Period,
		"StartDate":  req.StartDate,
		"EndDate":    req.EndDate,
	})
	if err != nil {
		log.Printf("Failed to create budget: %v", err)
		return &budgets.BudgetResponse{Succes: false, Message: "Failed to create budget"}, err
	}

	return &budgets.BudgetResponse{Succes: true, Message: "Budget created successfully"}, nil
}

func (s *AccountService) GetBudgetByid(ctx context.Context, req *budgets.GetBudgetByidRequest) (*budgets.GetBudgetByidResponse, error) {
	coll := s.db.Collection("budgets")

	var budget budgets.GetBudgetByidResponse
	err := coll.FindOne(ctx, bson.M{"id": req.Id}).Decode(&budget)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("No budget found with id: %v", req.Id)
			return nil, err
		}
		log.Printf("Failed to get budget by id: %v", err)
		return nil, err
	}

	return &budget, nil
}

func (s *AccountService) UpdateBudget(ctx context.Context, req *budgets.UpdateBudgetRequest) (*budgets.BudgetResponse, error) {
	coll := s.db.Collection("budgets")

	update := bson.M{}
	if req.Period != "" {
		update["period"] = req.Period
	}
	if req.StartDate != "" {
		update["StartDate"] = req.StartDate
	}
	if req.EndDate != "" {
		update["EndDate"] = req.EndDate
	}

	if len(update) == 0 {
		return &budgets.BudgetResponse{Succes: false, Message: "Nothing to update"}, nil
	}

	_, err := coll.UpdateOne(ctx, bson.M{"id": req.Id}, bson.M{"$set": update})
	if err != nil {
		log.Printf("Failed to update budget: %v", err)
		return &budgets.BudgetResponse{Succes: false, Message: "Failed to update budget"}, err
	}

	return &budgets.BudgetResponse{Succes: true, Message: "Budget updated successfully"}, nil
}

func (s *AccountService) DeleteBudget(ctx context.Context, req *budgets.DeleteBudgetRequest) (*budgets.BudgetResponse, error) {
	coll := s.db.Collection("budgets")

	_, err := coll.DeleteOne(ctx, bson.M{"id": req.Id})
	if err != nil {
		log.Printf("Failed to delete budget: %v", err)
		return &budgets.BudgetResponse{Succes: false, Message: "Failed to delete budget"}, err
	}

	return &budgets.BudgetResponse{Succes: true, Message: "Budget deleted successfully"}, nil
}

func (s *AccountService) ListBudgets(ctx context.Context, req *budgets.ListBudgetsRequest) (*budgets.ListBudgetsResponse, error) {
	coll := s.db.Collection("budgets")

	filter := bson.M{}
	if req.UserId != "" {
		filter["user_id"] = req.UserId
	}
	if req.CategoryId != "" {
		filter["category_id"] = req.CategoryId
	}
	if req.Amount != 0 {
		filter["amount"] = req.Amount
	}
	if req.Period != "" {
		filter["period"] = req.Period
	}
	if req.StartDate != "" {
		filter["start_date"] = req.StartDate
	}
	if req.EndDate != "" {
		filter["end_date"] = req.EndDate
	}

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		log.Printf("Failed to list budgets: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var budgetsList []*budgets.GetBudgetByidResponse
	for cursor.Next(ctx) {
		var budget budgets.GetBudgetByidResponse
		if err := cursor.Decode(&budget); err != nil {
			log.Printf("Failed to decode budget: %v", err)
			return nil, err
		}
		budgetsList = append(budgetsList, &budget)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}

	return &budgets.ListBudgetsResponse{Budgets: budgetsList}, nil
}

func (s *AccountService) UpdateBudgetAmount(ctx context.Context, UserId string, amount float32) error {
	coll := s.db.Collection("budgets")

	update := bson.M{
		"$inc": bson.M{
			"amount": -amount,
		},
	}
	_, err := coll.UpdateOne(ctx, bson.M{"UserId": UserId}, update)
	if err != nil {
		log.Printf("Failed to update account balance: %v", err)
		return err
	}
	return nil
}

func (s *AccountService) CheckBudget(ctx context.Context, userId string) (bool, error) {
	coll := s.db.Collection("budgets")

	// Define a struct to match the document structure
	var result struct {
		Amount    float32
		StartDate string
		EndDate   string
	}

	// Find the document for the given UserId
	err := coll.FindOne(ctx, bson.M{"UserId": userId}).Decode(&result)
	if err != nil {
		// Other errors (e.g., database issues)
		log.Printf("Failed to get budget by UserId: %v", err)
		return false, err
	}

	// Get the current date in the same string format as stored in MongoDB
	now := time.Now().Format("2006-01-02")

	// Check if 'now' is between 'StartDate' and 'EndDate'
	if now >= result.StartDate && now <= result.EndDate {

		// If within the date range, check if the amount is greater than 0
		if result.Amount <= 0 {

			return false, nil
		}
	}

	return true, nil
}
