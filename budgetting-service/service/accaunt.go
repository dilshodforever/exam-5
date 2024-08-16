package service

import (
	"context"
	"log"

	pb "github.com/dilshodforever/5-oyimtixon/genprotos/accaunts"
	s "github.com/dilshodforever/5-oyimtixon/storage"
)

type AccountService struct {
	stg s.InitRoot
	pb.UnimplementedAccountServiceServer
}

func NewAccountService(stg s.InitRoot) *AccountService {
	return &AccountService{stg: stg}
}

func (s *AccountService) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	resp, err := s.stg.Account().CreateAccount(ctx, req)
	if err != nil {
		log.Print(err)
	}
	
	return resp, err
}
						 	
func (s *AccountService) GetAccountByid(ctx context.Context, req *pb.GetByIdAccauntRequest) (*pb.GetAccountByidResponse, error) {
	resp, err := s.stg.Account().GetAccountByid(ctx, req)
	if err != nil {
		log.Print(err)
	}
	return resp, err
}

func (s *AccountService) UpdateAccount(ctx context.Context, req *pb.UpdateAccountRequest) (*pb.UpdateAccountResponse, error) {
	resp, err := s.stg.Account().UpdateAccount(ctx, req)
	if err != nil {
		log.Print(err)
	}
	return resp, err
}

func (s *AccountService) DeleteAccount(ctx context.Context, req *pb.DeleteAccountRequest) (*pb.UpdateAccountResponse, error) {
	resp, err := s.stg.Account().DeleteAccount(ctx, req)
	if err != nil {
		log.Print(err)
	}
	return resp, err
}

func (s *AccountService) ListAccounts(ctx context.Context, req *pb.ListAccountsRequest) (*pb.ListAccountsResponse, error) {
	resp, err := s.stg.Account().ListAccounts(ctx, req)
	if err != nil {
		log.Print(err)
	}
	return resp, err
}
