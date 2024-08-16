package postgres

import (
	pbAuth "github.com/dilshodforever/5-oyimtixon/genprotos/auth"
	pbUser "github.com/dilshodforever/5-oyimtixon/genprotos/user"
)

type InitRoot interface {
	Auth() Auth
	User() User
}

type Auth interface {
	Register(req *pbAuth.RegisterRequest) (*pbAuth.RegisterResponse, error)
	Login(req *pbAuth.LoginRequest) (*pbAuth.LoginResponse, error)
	ResetPassword(req *pbAuth.ResetPasswordRequest) (*pbAuth.ResetPasswordResponse, error)
}

type User interface {
	GetProfile(req *pbUser.GetProfileRequest) (*pbUser.GetProfileResponse, error)
	UpdateProfile(req *pbUser.UpdateProfileRequest) (*pbUser.UpdateProfileResponse, error)
	ChangePassword(req *pbUser.ChangePasswordRequest) (*pbUser.ChangePasswordResponse, error)
}
