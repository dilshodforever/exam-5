package service

import (
	"context"
	"log"

	pb "github.com/dilshodforever/5-oyimtixon/genprotos/goals"
	s "github.com/dilshodforever/5-oyimtixon/storage"
)

type GoalService struct {
	stg s.InitRoot
	pb.UnimplementedGoalServiceServer
}

func NewGoalService(stg s.InitRoot) *GoalService {
	return &GoalService{stg: stg}
}

func (s *GoalService) CreateGoal(ctx context.Context, req *pb.CreateGoalRequest) (*pb.GoalResponse, error) {
	resp, err := s.stg.Goal().CreateGoal(ctx, req)
	if err != nil {
		log.Print(err)
	}
	return resp, err
}

func (s *GoalService) GetGoalByid(ctx context.Context, req *pb.GetGoalByidRequest) (*pb.GetGoalResponse, error) {
	resp, err := s.stg.Goal().GetGoalByid(ctx, req)
	if err != nil {
		log.Print(err)
	}
	return resp, err
}

func (s *GoalService) UpdateGoal(ctx context.Context, req *pb.UpdateGoalRequest) (*pb.GoalResponse, error) {
	resp, err := s.stg.Goal().UpdateGoal(ctx, req)
	if err != nil {
		log.Print(err)
	}
	return resp, err
}

func (s *GoalService) DeleteGoal(ctx context.Context, req *pb.DeleteGoalRequest) (*pb.GoalResponse, error) {
	resp, err := s.stg.Goal().DeleteGoal(ctx, req)
	if err != nil {
		log.Print(err)
	}
	return resp, err
}

func (s *GoalService) ListGoals(ctx context.Context, req *pb.ListGoalsRequest) (*pb.ListGoalsResponse, error) {
	resp, err := s.stg.Goal().ListGoals(ctx, req)
	if err != nil {
		log.Print(err)
	}
	return resp, err
}
