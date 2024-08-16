package service

import (
	"context"
	"log"

	pb "github.com/dilshodforever/5-oyimtixon/genprotos/categories"
	s "github.com/dilshodforever/5-oyimtixon/storage"
)

type CategoryService struct {
	stg s.InitRoot
	pb.UnimplementedCategoryServiceServer
}

func NewCategoryService(stg s.InitRoot) *CategoryService {
	return &CategoryService{stg: stg}
}

func (s *CategoryService) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.CategoryResponse, error) {
	resp, err := s.stg.Category().CreateCategory(ctx, req)
	if err != nil {
		log.Print(err)
	}
	return resp, err
}

func (s *CategoryService) UpdateCategory(ctx context.Context, req *pb.UpdateCategoryRequest) (*pb.CategoryResponse, error) {
	resp, err := s.stg.Category().UpdateCategory(ctx, req)
	if err != nil {
		log.Print(err)
	}
	return resp, err
}

func (s *CategoryService) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryRequest) (*pb.CategoryResponse, error) {
	resp, err := s.stg.Category().DeleteCategory(ctx, req)
	if err != nil {
		log.Print(err)
	}
	return resp, err
}

func (s *CategoryService) ListCategories(ctx context.Context, req *pb.ListCategoriesRequest) (*pb.ListCategoriesResponse, error) {
	resp, err := s.stg.Category().ListCategories(ctx, req)
	if err != nil {
		log.Print(err)
	}
	return resp, err
}

func (s *CategoryService) GetByidCategory(ctx context.Context, req *pb.GetByidCategoriesRequest) (*pb.GetByidCategoriesResponse, error) {
	resp, err := s.stg.Category().GetByidCategory(ctx, req)
	if err != nil {
		log.Print(err)
	}
	return resp, err
}
