package service

import (
	"context"
	"log"

	pb "github.com/dilshodforever/5-oyimtixon/genprotos/budgets"
	s "github.com/dilshodforever/5-oyimtixon/storage"
)

type BudgetService struct {
	stg s.InitRoot
	pb.UnimplementedBudgetServiceServer
}

func NewBudgetService(stg s.InitRoot) *BudgetService {
	return &BudgetService{stg: stg}
}

func (s *BudgetService) CreateBudget(ctx context.Context, req *pb.CreateBudgetRequest) (*pb.BudgetResponse, error) {
	resp, err := s.stg.Budget().CreateBudget(ctx, req)
	if err != nil {
		log.Print(err)
	}
	return resp, err
}

func (s *BudgetService) GetBudgetByid(ctx context.Context, req *pb.GetBudgetByidRequest) (*pb.GetBudgetByidResponse, error) {
	resp, err := s.stg.Budget().GetBudgetByid(ctx, req)
	if err != nil {
		log.Print(err)
	}
	return resp, err
}

func (s *BudgetService) UpdateBudget(ctx context.Context, req *pb.UpdateBudgetRequest) (*pb.BudgetResponse, error) {
	resp, err := s.stg.Budget().UpdateBudget(ctx, req)
	if err != nil {
		log.Print(err)
	}
	return resp, err
}

func (s *BudgetService) DeleteBudget(ctx context.Context, req *pb.DeleteBudgetRequest) (*pb.BudgetResponse, error) {
	resp, err := s.stg.Budget().DeleteBudget(ctx, req)
	if err != nil {
		log.Print(err)
	}
	return resp, err
}

func (s *BudgetService) ListBudgets(ctx context.Context, req *pb.ListBudgetsRequest) (*pb.ListBudgetsResponse, error) {
	resp, err := s.stg.Budget().ListBudgets(ctx, req)
	if err != nil {
		log.Print(err)
	}
	return resp, err
}
