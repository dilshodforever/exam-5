package postgres

import (
	"database/sql"

	pb "github.com/dilshodforever/5-oyimtixon/genprotos/user"
)

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{db: db}
}

func (p *UserStorage) GetProfile(req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	var profile pb.GetProfileResponse
	query := `
		SELECT id, username, email, phone_number
		FROM users
		WHERE id = $1
	`
	err := p.db.QueryRow(query, req.Id).Scan(
		&profile.Id, &profile.Username, &profile.Email,
		&profile.PhoneNumber,
	)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (p *UserStorage) UpdateProfile(req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	query := `
		UPDATE users
		SET phone_number = $1, username=$2, email=&3 updated_at = now()
		WHERE username = $3
		RETURNING id, username, email, phone_number, updated_at
	`
	var profile pb.UpdateProfileResponse
	err := p.db.QueryRow(query, req.PhoneNumber,req.Username, req.Email).Scan(
		&profile.Id, &profile.Username, &profile.Email,
		&profile.PhoneNumber, &profile.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (p *UserStorage) ChangePassword(req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	query := `
		UPDATE users
		SET password_hash = $1, updated_at = now()
		WHERE id = $2
	`
	_, err := p.db.Exec(query, req.NewPassword, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.ChangePasswordResponse{Message: "Password successfully changed"}, nil
}
