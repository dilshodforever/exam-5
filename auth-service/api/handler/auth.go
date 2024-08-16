package handler

import (
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"time"

	"github.com/dilshodforever/5-oyimtixon/api/token"
	"github.com/go-redis/redis"

	pb "github.com/dilshodforever/5-oyimtixon/genprotos/auth"

	"github.com/gin-gonic/gin"
)

// Register handles the creation of a new user
// @Summary Register a new user
// @Description Register a new user with username, email, password, and full name
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Register body pb.Register true "Register user"
// @Success 200 {object} pb.RegisterResponse "Registration successful"
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal server error"
// @Router /auth/register [post]
func (h *Handler) Register(ctx *gin.Context) {
	req := &pb.RegisterRequest{}
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	res, err := h.Auth.Register(ctx, req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		slog.Info(err.Error())
		return
	}
	tokens := token.GenereteJWTToken(res)
	err = h.Redis.SaveToken(res.Id, tokens.RefreshToken, 30*24*time.Hour)
	if err != nil {
		slog.Info(err.Error())
	}
	ctx.JSON(200, res)
}

// Login handles user login
// @Summary Login a user
// @Description Login a user with username and password
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Login body pb.LoginRequest true "Login user"
// @Success 200 {object} pb.LoginResponse "Login successful"
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/login [post]
func (h *Handler) Login(ctx *gin.Context) {
	req := &pb.LoginRequest{}
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	res, err := h.Auth.Login(ctx, req)
	if err != nil {
		ctx.JSON(401, gin.H{"error": err})
		return
	}
	if !res.Success {
		ctx.JSON(400, res)
		return
	}
	fmt.Println(res.Id)
	token, err := h.Redis.Get("token:"+res.Id)
	if err != nil {
		if err == redis.Nil {
			log.Printf("Key %s does not exist", res.Id)
		} else {
			log.Printf("Error getting key %s: %v", res.Id, err)
		}
	}
	ctx.JSON(200, token)
}

// ForgotPassword handles forgot password requests
// @Summary Forgot password
// @Description Request password reset
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param ForgotPassword body pb.ForgotPasswordRequest true "Forgot password"
// @Success 200 {object} pb.ForgotPasswordResponse "Password reset instructions sent"
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal server error"
// @Router /auth/forgot-password [post]
func (h *Handler) ForgotPassword(ctx *gin.Context) {
	req := &pb.ForgotPasswordRequest{}
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	err = h.Redis.SaveEmailCode(req.Email, code, 10*time.Minute)
	if err != nil {
		slog.Info(err.Error())
	}
	ctx.JSON(200, "Code: "+code)
}

// ResetPassword handles password reset
// @Summary Reset password
// @Description Reset password with token and new password
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param ResetPassword body pb.ResetPassword true "Reset password"
// @Success 200 {object} pb.ResetPasswordResponse "Password successfully reset"
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal server error"
// @Router /auth/reset-password [put]
func (h *Handler) ResetPassword(ctx *gin.Context) {
	req := &pb.ResetPasswordRequest{}
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	jwtToken := ctx.Request.Header.Get("Authorization")
	claims, err := token.ExtractClaim(jwtToken)
	if err != nil {

		panic(err)
	}
	req.Email = claims["email"].(string)
	req.Username = claims["username"].(string)
	code, err := h.Redis.Get("email_code:" + req.Email)
	if err != nil {
		slog.Info(err.Error())
	}

	if code != req.EmailPassword {
		ctx.JSON(400, "error while resetting password. Please check your email and try again")
		return
	}

	res, err := h.Auth.ResetPassword(ctx, req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	ctx.JSON(200, res)
}
