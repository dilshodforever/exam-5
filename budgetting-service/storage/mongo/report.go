package postgres

import (
	"context"
	"log"

	"github.com/dilshodforever/5-oyimtixon/genprotos/reports"
	"go.mongodb.org/mongo-driver/bson"
)
func (s *AccountService) GetReports(ctx context.Context, req *reports.GenerateReportRequest) (*reports.GenerateReportResponse, error) {
	coll := s.db.Collection("transactions")
	filter := bson.M{
		"UserId": req.UserId,
		"type":   req.Type,
	}

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		log.Printf("Failed to list transactions: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var responses reports.GenerateReportResponse
	for cursor.Next(ctx) {
		var response reports.GetSpending
		if err := cursor.Decode(&response); err != nil {
			log.Printf("Failed to decode transaction: %v", err)
			return nil, err
		}
		responses.Reeports = append(responses.Reeports, &response)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}

	return &responses, nil
}
