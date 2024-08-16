package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dilshodforever/5-oyimtixon/genprotos/goals"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *AccountService) CreateGoal(ctx context.Context, req *goals.CreateGoalRequest) (*goals.GoalResponse, error) {
	coll := s.db.Collection("goals")
	id := uuid.NewString()
	_, err := coll.InsertOne(ctx, bson.M{
		"id":            id,
		"UserId":        req.UserId,
		"name":          req.Name,
		"TargetAmount":  req.TargetAmount,
		"CurrentAmount": 0,
		"deadline":      req.Deadline,
		"status":        "Inprogres",
	})
	if err != nil {
		log.Printf("Failed to create goal: %v", err)
		return &goals.GoalResponse{Success: false, Message: "Failed to create goal"}, err
	}

	return &goals.GoalResponse{Success: true, Message: "Goal created successfully"}, nil
}

func (s *AccountService) GetGoalByid(ctx context.Context, req *goals.GetGoalByidRequest) (*goals.GetGoalResponse, error) {
	coll := s.db.Collection("goals")

	var goal goals.GetGoalResponse
	err := coll.FindOne(ctx, bson.M{"id": req.Id}).Decode(&goal)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("No goal found with id: %v", req.Id)
			return nil, err
		}
		log.Printf("Failed to get goal by id: %v", err)
		return nil, err
	}

	return &goal, nil
}

func (s *AccountService) UpdateGoal(ctx context.Context, req *goals.UpdateGoalRequest) (*goals.GoalResponse, error) {
	coll := s.db.Collection("goals")

	update := bson.M{}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.TargetAmount > 0 {
		update["target_amount"] = req.TargetAmount
	}
	if req.CurrentAmount > 0 {
		update["current_amount"] = req.CurrentAmount
	}
	if req.Deadline != "" {
		update["deadline"] = req.Deadline
	}
	if req.Status != "" {
		update["status"] = req.Status
	}
	fmt.Println("filter", update)
	if len(update) == 0 {
		return &goals.GoalResponse{Success: false, Message: "Nothing to update"}, nil
	}

	_, err := coll.UpdateOne(ctx, bson.M{"id": req.Id}, bson.M{"$set": update})
	if err != nil {
		log.Printf("Failed to update goal: %v", err)
		return &goals.GoalResponse{Success: false, Message: "Failed to update goal"}, err
	}

	return &goals.GoalResponse{Success: true, Message: "Goal updated successfully"}, nil
}

func (s *AccountService) DeleteGoal(ctx context.Context, req *goals.DeleteGoalRequest) (*goals.GoalResponse, error) {
	coll := s.db.Collection("goals")

	_, err := coll.DeleteOne(ctx, bson.M{"id": req.Id})
	if err != nil {
		log.Printf("Failed to delete goal: %v", err)
		return &goals.GoalResponse{Success: false, Message: "Failed to delete goal"}, err
	}

	return &goals.GoalResponse{Success: true, Message: "Goal deleted successfully"}, nil
}

func (s *AccountService) ListGoals(ctx context.Context, req *goals.ListGoalsRequest) (*goals.ListGoalsResponse, error) {
	coll := s.db.Collection("goals")

	filter := bson.M{}
	if req.UserId != "" {
		filter["user_id"] = req.UserId
	}
	if req.Name != "" {
		filter["name"] = req.Name
	}
	if req.TargetAmount > 0 {
		filter["target_amount"] = req.TargetAmount
	}
	if req.CurrentAmount > 0 {
		filter["current_amount"] = req.CurrentAmount
	}
	if req.Deadline != "" {
		
		filter["deadline"] = req.Deadline
	}
	if req.Status != "" {
		
		filter["status"] = req.Status
	}
	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		log.Printf("Failed to list goals: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var goalsList []*goals.GetGoalResponse
	for cursor.Next(ctx) {
		var goal goals.GetGoalResponse
		if err := cursor.Decode(&goal); err != nil {
			log.Printf("Failed to decode goal: %v", err)
			return nil, err
		}
		goalsList = append(goalsList, &goal)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}

	return &goals.ListGoalsResponse{Goals: goalsList}, nil
}

func (s *AccountService) UpdateGoulAmount(ctx context.Context, UserId string, amount float32) error {
	coll := s.db.Collection("goals")

	update := bson.M{
		"$inc": bson.M{
			"CurrentAmount": +amount,
		},
	}
	_, err := coll.UpdateOne(ctx, bson.M{"UserId": UserId}, update)
	if err != nil {
		log.Printf("Failed to update account balance: %v", err)
		return err
	}
	return nil
}

func (s *AccountService) CheckGoal(ctx context.Context, userId string) (bool, string, error) {
	coll := s.db.Collection("goals")

	// Define a struct to match the document structure
	var result struct {
		TargetAmount  float32
		CurrentAmount float32
		Deadline      string
	}

	// Find the document for the given UserId
	err := coll.FindOne(ctx, bson.M{"UserId": userId}).Decode(&result)
	if err != nil {
		log.Printf("Failed to get goal by UserId: %v", err)
		return false, "", err
	}

	// Get the current date in the same string format as stored in MongoDB
	now := time.Now().Format("2006-01-02")

	// Check if 'now' matches the 'Deadline'
	if now == result.Deadline {
		if result.CurrentAmount < result.TargetAmount {
			err= s.UpdateStatusByUserId(ctx,userId,"Filed")
			if err!=nil{
				log.Print("Error while update goal status")
				return false, "", err
			}
			return false, "The deadline has passed, and you did not reach your savings goal.", nil
		}
		   err= s.UpdateStatusByUserId(ctx,userId,"Success")
			if err!=nil{
				log.Print("Error while update goal status")
				return false, "", err
			}
		return false, "Congratulations! You reached your savings goal by the deadline.", nil
	}

	return true, "", nil
}




func (s *AccountService) UpdateStatusByUserId(ctx context.Context, userid string, status string) (error) {
	coll := s.db.Collection("goals")
	update := bson.M{
		"$set": bson.M{
			"status":     status,
		},
	}
	_, err := coll.UpdateOne(ctx, bson.M{"UserId": userid}, update)
	if err != nil {
		log.Printf("Failed to update account: %v", err)
		return err
	}
	return nil
}
