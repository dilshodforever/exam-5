package postgres

import (
	"database/sql"
	"fmt"

	"github.com/dilshodforever/5-oyimtixon/config"
	u "github.com/dilshodforever/5-oyimtixon/storage"
	"github.com/go-redis/redis/v8"

	_ "github.com/lib/pq"
)

type Storage struct {
	Db    *sql.DB
	Users u.User
	Auths u.Auth
	rdb   redis.Client
}

func NewPostgresStorage() (u.InitRoot, error) {
	config := config.Load()
	con := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=require",
		config.PostgresUser, config.PostgresPassword,
		config.PostgresHost, config.PostgresPort,
		config.PostgresDatabase)
	db, err := sql.Open("postgres", con)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
	
		return nil, err
	}
	fmt.Println("Connection to database!!!")
	return &Storage{Db: db}, err

}

func (s *Storage) User() u.User {
	if s.Users == nil {
		s.Users = &UserStorage{s.Db}
	}
	return s.Users
}

func (s *Storage) Auth() u.Auth {
	if s.Auths == nil {
		s.Auths = &AuthStorage{db: s.Db, client: &s.rdb}
	}
	return s.Auths
}