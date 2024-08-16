package handler

import (
	"github.com/dilshodforever/5-oyimtixon/service"
)

type Handler struct {
	Auth  *service.AuthService
	User  *service.UserService
	Redis InMemoryStorageI
}

func NewHandler(auth *service.AuthService, user *service.UserService,
	redis InMemoryStorageI) *Handler {
	return &Handler{
		Auth:  auth,
		User:  user,
		Redis: redis,
	}

}
