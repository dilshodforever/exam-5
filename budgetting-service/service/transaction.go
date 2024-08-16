package service

import (
	"context"
	"fmt"
	"log"

	pb "github.com/dilshodforever/5-oyimtixon/genprotos/transactions"
	"github.com/dilshodforever/5-oyimtixon/kafkaconnect"
	sender "github.com/dilshodforever/5-oyimtixon/kafkasender"
	"github.com/dilshodforever/5-oyimtixon/model"
	s "github.com/dilshodforever/5-oyimtixon/storage"
)

type TransactionService struct {
	stg s.InitRoot
	pb.UnimplementedTransactionServiceServer
}

func NewTransactionService(stg s.InitRoot) *TransactionService {
	return &TransactionService{stg: stg}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.TransactionResponse, error) {
	resp, err := s.stg.Transaction().CreateTransaction(ctx, req)
	if err != nil {
		log.Printf("Failed to create transaction: %v", err)
		return &pb.TransactionResponse{Success: false, Message: "Failed to create transaction"}, err
	}
	if req.Type=="-"{
	err = s.stg.Account().UpdateBalanceMinus(ctx, req.AccountId, req.Amount)
	if err != nil {
		log.Printf("Failed to update balance: %v", err)
		return &pb.TransactionResponse{Success: false, Message: "Failed to update balance"}, err
	}

	err = s.stg.Budget().UpdateBudgetAmount(ctx, req.UserId, req.Amount)
	if err != nil {
		log.Printf("Failed to update budjet balance: %v", err)
		return &pb.TransactionResponse{Success: false, Message: "Failed to update balance"}, err
	}
	check, err := s.stg.Account().CheckBudget(ctx, req.UserId)
	if err != nil {
		log.Printf("Failed to check budjet balance: %v", err)
		return &pb.TransactionResponse{Success: false, Message: "Failed to update balance"}, err
	}
	if !check {
		kaf:=kafkaconnect.ConnectToKafka()
		request := model.Send{Message: "Your Budget was empty", Userid: req.UserId}
		err=sender.CreateNotification(kaf,request)
		if err != nil {
			fmt.Println("error while send kafka")
			return &pb.TransactionResponse{Success: false, Message: "Failed send kafka"}, err
		}
	}
} else{
	err = s.stg.Account().UpdateBalance(ctx, req.AccountId, req.Amount)
	if err != nil {
		log.Printf("Failed to update balance: %v", err)
		return &pb.TransactionResponse{Success: false, Message: "Failed to update balance"}, err
	}
	err = s.stg.Goal().UpdateGoulAmount(ctx, req.UserId, req.Amount)
	if err != nil {
		log.Printf("Failed to update goal balance: %v", err)
		return &pb.TransactionResponse{Success: false, Message: "Failed to update goal balance"}, err
	}

	goalcheck, message, err:=s.stg.Goal().CheckGoal(ctx, req.UserId)
	if err != nil {
		log.Printf("Failed to cheak goal balance: %v", err)
		return &pb.TransactionResponse{Success: false, Message: "Failed to update goal balance"}, err
	}
	if !goalcheck{
		kaf:=kafkaconnect.ConnectToKafka()
		request := model.Send{Message: message, Userid: req.UserId}
		err=sender.CreateNotification(kaf,request)
		if err != nil {
			fmt.Println("error while send kafka")
			return &pb.TransactionResponse{Success: false, Message: "Failed send kafka"}, err
		}
	}
	
}
	
	return resp, nil
}

func (s *TransactionService) GetTransaction(ctx context.Context, req *pb.GetTransactionRequest) (*pb.GetTransactionResponse, error) {
	resp, err := s.stg.Transaction().GetTransaction(ctx, req)
	if err != nil {
		log.Printf("Failed to get transaction: %v", err)
		return nil, err
	}
	return resp, nil
}

func (s *TransactionService) UpdateTransaction(ctx context.Context, req *pb.UpdateTransactionRequest) (*pb.TransactionResponse, error) {
	resp, err := s.stg.Transaction().UpdateTransaction(ctx, req)
	if err != nil {
		log.Printf("Failed to update transaction: %v", err)
		return &pb.TransactionResponse{Success: false, Message: "Failed to update transaction"}, err
	}
	return resp, nil
}

func (s *TransactionService) DeleteTransaction(ctx context.Context, req *pb.DeleteTransactionRequest) (*pb.TransactionResponse, error) {
	resp, err := s.stg.Transaction().DeleteTransaction(ctx, req)
	if err != nil {
		log.Printf("Failed to delete transaction: %v", err)
		return &pb.TransactionResponse{Success: false, Message: "Failed to delete transaction"}, err
	}
	return resp, nil
}

func (s *TransactionService) ListTransactions(ctx context.Context, req *pb.ListTransactionsRequest) (*pb.ListTransactionsResponse, error) {
	resp, err := s.stg.Transaction().ListTransactions(ctx, req)
	if err != nil {
		log.Printf("Failed to list transactions: %v", err)
		return nil, err
	}
	return resp, nil
}
