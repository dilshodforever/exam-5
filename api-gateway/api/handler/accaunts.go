package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/dilshodforever/5-oyimtixon/api/middleware"
	pb "github.com/dilshodforever/5-oyimtixon/genprotos/accaunts"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// CreateAccount handles creating a new account
// @Summary      Create Account
// @Description  Create a new account
// @Tags         Account
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        account body pb.CreateAccountRequest true "Account details"
// @Success      200 {object} pb.CreateAccountResponse "Account created successfully"
// @Failure      400 {string} string "Invalid input"
// @Failure      500 {string} string "Error while creating account"
// @Router       /account/create [post]
func (h *Handler) CreateAccount(ctx *gin.Context) {
	var req pb.CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	id := middleware.GetUserId(ctx)
	req.UserId = id
	res, err := h.Account.CreateAccount(ctx, &req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}

// GetAccountById handles retrieving an account by ID
// @Summary      Get Account by ID
// @Description  Get details of an account by ID
// @Tags         Account
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id query string true "Account ID"
// @Success      200 {object} pb.GetAccountByidResponse "Account details"
// @Failure      400 {string} string "Missing or invalid ID"
// @Failure      500 {string} string "Error while fetching account"
// @Router       /account/get/{id} [get]
func (h *Handler) GetAccountById(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		ctx.JSON(400, gin.H{"error": "Missing or invalid ID"})
		return
	}

	req := &pb.GetByIdAccauntRequest{Id: id}

	// Try to get the account data from Redis
	cacheKey := "user_id:" + req.Id
	cachedData, err := h.Redis.Get(cacheKey)
	if err != nil && err != redis.Nil {
		slog.Error("Failed to get data from Redis: %v", err)
	}

	var res *pb.GetAccountByidResponse

	if err == nil {
		// Data found in Redis, unmarshal it
		fmt.Println("redis")
		if err := json.Unmarshal([]byte(cachedData), &res); err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to parse cached data"})
			return
		}
	} else if err == redis.Nil {
		// Data not found in Redis, fetch it from the service
		fmt.Println("mongo")
		res, err = h.Account.GetAccountByid(ctx, req)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		data, err := json.Marshal(res)
		if err != nil {
			slog.Error("Failed to marshal response: %v", err)
		} else {
			// Convert the byte array to a string
			dataStr := string(data)

			// Store the data in Redis with a 30-minute expiration
			err := h.Redis.Set(cacheKey, dataStr, 30*time.Minute)
			if err != nil {
				slog.Error("Failed to set data in Redis: %v", err)
			}
		}
	} else {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}

// UpdateAccount handles updating an account
// @Summary      Update Account
// @Description  Update details of an account
// @Tags         Account
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        account body pb.UpdateAccountRequest true "Updated account details"
// @Success      200 {object} pb.UpdateAccountResponse "Account updated successfully"
// @Failure      400 {string} string "Invalid input"
// @Failure      500 {string} string "Error while updating account"
// @Router       /account/update [put]
func (h *Handler) UpdateAccount(ctx *gin.Context) {
	var req pb.UpdateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	id := middleware.GetUserId(ctx)
	req.Id = id
	res, err := h.Account.UpdateAccount(ctx, &req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}

// DeleteAccount handles deleting an account
// @Summary      Delete Account
// @Description  Delete an account by ID
// @Tags         Account
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} pb.UpdateAccountResponse "Account deleted successfully"
// @Failure      400 {string} string "Missing or invalid ID"
// @Failure      500 {string} string "Error while deleting account"
// @Router       /account/delete [delete]
func (h *Handler) DeleteAccount(ctx *gin.Context) {
	id := middleware.GetUserId(ctx)
	req := &pb.DeleteAccountRequest{Id: id}
	res, err := h.Account.DeleteAccount(ctx, req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}

// ListAccounts handles listing all accounts
// @Summary      List Accounts
// @Description  Get a list of all accounts based on the provided query parameters
// @Tags         Account
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        name     query     string  false  "Name of the account"  example("Savings")
// @Param        type     query     string  false  "Type of the account"  example("Checking")
// @Param        balance  query     number   false  "Balance of the account"  example(1000.50)
// @Param        currency query     string  false  "Currency of the account"  example("USD")
// @Success      200      {object}  pb.ListAccountsResponse "List of accounts"
// @Failure      500      {string}  string                  "Error while fetching accounts"
// @Router       /account/list [get]
func (h *Handler) ListAccounts(ctx *gin.Context) {
	req := &pb.ListAccountsRequest{
		Name:     ctx.Query("name"),
		Type:     ctx.Query("type"),
		Currency: ctx.Query("currency"),
	}
	// Convert balance to float
	if balanceStr := ctx.Query("balance"); balanceStr != "" {
		if balance, err := strconv.ParseFloat(balanceStr, 32); err == nil {
			req.Balance = float32(balance)
		}
	}

	res, err := h.Account.ListAccounts(ctx, req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}
