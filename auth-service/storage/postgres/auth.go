package postgres

import (
	"database/sql"
	"log/slog"

	pb "github.com/dilshodforever/5-oyimtixon/genprotos/auth"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type AuthStorage struct {
	db     *sql.DB
	client *redis.Client
}

func NewAuthStorage(db *sql.DB) *AuthStorage {
	return &AuthStorage{db: db}
}

func (p *AuthStorage) Register(req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	userId := uuid.NewString()
	query := `
		INSERT INTO users (id, username, email, password_hash, phone_number, role)
		VALUES ($1, $2, $3, $4, $5, 'user')
		RETURNING id, username, email, phone_number, created_at, role
	`
	var user pb.RegisterResponse
	err := p.db.QueryRow(query, userId, req.Name, req.Email, req.Password, req.PhoneNumber).Scan(
		&user.Id, &user.Name, &user.Email, &user.PhoneNumber, &user.CreatedAt, &user.Role,
	)
	if err != nil {
		slog.Info(err.Error())
		return nil, err
	}
	return &user, nil
}

func (p *AuthStorage) Login(req *pb.LoginRequest) (*pb.LoginResponse, error) {
	query := `
		SELECT id
		FROM users
		WHERE username = $1 AND password_hash = $2
	`
	var id string
	err := p.db.QueryRow(query, req.Name, req.Password).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.LoginResponse{Message: "Invalid username or password", Success: false}, nil
		}
		return nil, err
	}
	
	return &pb.LoginResponse{Message: "Login successful", Success: true, Id: id}, nil
}

func (p *AuthStorage) ResetPassword(req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	query := `
		UPDATE users
		SET password_hash = $1
		WHERE email = $2 and username=$3
	`
	_, err := p.db.Exec(query, req.NewPassword, req.Email, req.Username)
	if err != nil {
		return nil, err
	}
	return &pb.ResetPasswordResponse{Message: "Password reset successfully"}, nil
}
