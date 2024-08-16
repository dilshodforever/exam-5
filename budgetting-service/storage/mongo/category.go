package postgres

import (
	"context"
	"log"

	"github.com/dilshodforever/5-oyimtixon/genprotos/categories"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *AccountService) CreateCategory(ctx context.Context, req *categories.CreateCategoryRequest) (*categories.CategoryResponse, error) {
	coll := s.db.Collection("categories")
	id := uuid.NewString()
	_, err := coll.InsertOne(ctx, bson.M{
		"id":     id,
		"UserId": req.UserId,
		"name":   req.Name,
		"type":   req.Type,
	})
	if err != nil {
		log.Printf("Failed to create category: %v", err)
		return &categories.CategoryResponse{Success: false, Message: "Failed to create category"}, err
	}

	return &categories.CategoryResponse{Success: true, Message: "Category created successfully"}, nil
}

func (s *AccountService) UpdateCategory(ctx context.Context, req *categories.UpdateCategoryRequest) (*categories.CategoryResponse, error) {
	coll := s.db.Collection("categories")

	update := bson.M{}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Type != "" {
		update["type"] = req.Type
	}

	if len(update) == 0 {
		return &categories.CategoryResponse{Success: false, Message: "Nothing to update"}, nil
	}

	_, err := coll.UpdateOne(ctx, bson.M{"id": req.Id}, bson.M{"$set": update})
	if err != nil {
		log.Printf("Failed to update category: %v", err)
		return &categories.CategoryResponse{Success: false, Message: "Failed to update category"}, err
	}

	return &categories.CategoryResponse{Success: true, Message: "Category updated successfully"}, nil
}

func (s *AccountService) DeleteCategory(ctx context.Context, req *categories.DeleteCategoryRequest) (*categories.CategoryResponse, error) {
	coll := s.db.Collection("categories")

	_, err := coll.DeleteOne(ctx, bson.M{"id": req.Id})
	if err != nil {
		log.Printf("Failed to delete category: %v", err)
		return &categories.CategoryResponse{Success: false, Message: "Failed to delete category"}, err
	}

	return &categories.CategoryResponse{Success: true, Message: "Category deleted successfully"}, nil
}

func (s *AccountService) ListCategories(ctx context.Context, req *categories.ListCategoriesRequest) (*categories.ListCategoriesResponse, error) {
	coll := s.db.Collection("categories")

	filter := bson.M{}
	if req.UserId != "" {
		filter["user_id"] = req.UserId
	}
	if req.Name != "" {
		filter["name"] = req.Name
	}
	if req.Type != "" {
		filter["type"] = req.Type
	}

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		log.Printf("Failed to list categories: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var categoriesList []*categories.GetByidCategoriesResponse
	for cursor.Next(ctx) {
		var category categories.GetByidCategoriesResponse
		if err := cursor.Decode(&category); err != nil {
			log.Printf("Failed to decode category: %v", err)
			return nil, err
		}
		categoriesList = append(categoriesList, &category)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}

	return &categories.ListCategoriesResponse{Categories: categoriesList}, nil
}

func (s *AccountService) GetByidCategory(ctx context.Context, req *categories.GetByidCategoriesRequest) (*categories.GetByidCategoriesResponse, error) {
	coll := s.db.Collection("categories")

	var category categories.GetByidCategoriesResponse
	err := coll.FindOne(ctx, bson.M{"_id": req.Id}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("No category found with id: %v", req.Id)
			return nil, err
		}
		log.Printf("Failed to get category by id: %v", err)
		return nil, err
	}

	return &category, nil
}
