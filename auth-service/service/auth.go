package service

import (
	"context"
	"log"

	pb "github.com/dilshodforever/5-oyimtixon/genprotos/auth"
	s "github.com/dilshodforever/5-oyimtixon/storage"
)

type AuthService struct {
	stg s.InitRoot
	pb.UnimplementedAuthServiceServer
}

func NewAuthService(stg s.InitRoot) *AuthService {
	return &AuthService{stg: stg}
}

func (a *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	res, err := a.stg.Auth().Register(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return res, nil
}

func (a *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	res, err := a.stg.Auth().Login(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return res, nil
}

func (a *AuthService) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	res, err := a.stg.Auth().ResetPassword(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return res, nil
}