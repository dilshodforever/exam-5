package handler

import (
	"strconv"
	"time"

	"github.com/dilshodforever/5-oyimtixon/api/middleware"
	pb "github.com/dilshodforever/5-oyimtixon/genprotos/transactions"
	"github.com/dilshodforever/5-oyimtixon/kafkasender"
	"github.com/eapache/go-resiliency/breaker"
	"github.com/gin-gonic/gin"
)

var (
	transactionBreaker = breaker.New(3, 1, time.Minute)
)

// CreateTransaction handles creating a new transaction
// @Summary      Create Transaction
// @Description  Create a new transaction
// @Tags         Transaction
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        transaction body pb.CreateTransactionRequest true "Transaction details"
// @Success      200 {object} string "Transaction created successfully"
// @Failure      400 {string} string "Invalid input"
// @Failure      500 {string} string "Error while creating transaction"
// @Router       /transaction/create [post]
func (h *Handler) CreateTransaction(ctx *gin.Context) {
	var req pb.CreateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	id := middleware.GetUserId(ctx)
	req.UserId = id
	err := transactionBreaker.Run(func() error {
		_, err := kafkasender.CreateTransaction(h.Kafka, &req)
		if err != nil {
			return err
		}

		ctx.JSON(200, "Success")
		return nil
	})

	if err != nil {
		if err == breaker.ErrBreakerOpen {
			ctx.JSON(500, gin.H{"error": "Service is temporarily unavailable. Please try again later."})
		} else {
			ctx.JSON(500, gin.H{"error": err.Error()})
		}
		return
	}
}

// GetTransaction handles retrieving a transaction by ID
// @Summary      Get Transaction by ID
// @Description  Retrieve details of a transaction by ID
// @Tags         Transaction
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Transaction ID"
// @Success      200 {object} pb.GetTransactionResponse "Transaction details"
// @Failure      400 {string} string "Missing or invalid ID"
// @Failure      500 {string} string "Error while fetching transaction"
// @Router       /transaction/get/{id} [get]
func (h *Handler) GetTransaction(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(400, gin.H{"error": "Missing or invalid ID"})
		return
	}

	req := &pb.GetTransactionRequest{Id: id}

	res, err := h.Transaction.GetTransaction(ctx, req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}

// UpdateTransaction handles updating a transaction
// @Summary      Update Transaction
// @Description  Update details of a transaction
// @Tags         Transaction
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        transaction body pb.UpdateTransactionRequest true "Updated transaction details"
// @Success      200 {object} pb.TransactionResponse "Transaction updated successfully"
// @Failure      400 {string} string "Invalid input"
// @Failure      500 {string} string "Error while updating transaction"
// @Router       /transaction/update [put]
func (h *Handler) UpdateTransaction(ctx *gin.Context) {
	var req pb.UpdateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	res, err := h.Transaction.UpdateTransaction(ctx, &req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}

// DeleteTransaction handles deleting a transaction
// @Summary      Delete Transaction
// @Description  Delete a transaction by ID
// @Tags         Transaction
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Transaction ID"
// @Success      200 {object} pb.TransactionResponse "Transaction deleted successfully"
// @Failure      400 {string} string "Missing or invalid ID"
// @Failure      500 {string} string "Error while deleting transaction"
// @Router       /transaction/delete/{id} [delete]
func (h *Handler) DeleteTransaction(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(400, gin.H{"error": "Missing or invalid ID"})
		return
	}

	req := &pb.DeleteTransactionRequest{Id: id}

	res, err := h.Transaction.DeleteTransaction(ctx, req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}

// ListTransactions handles listing all transactions
// @Summary      List Transactions
// @Description  Get a list of all transactions based on query parameters
// @Tags         Transaction
// @Accept       json
// @Produce      json
// @Param        account_id    query  string  false  "Account ID"
// @Param        category_id   query  string  false  "Category ID"
// @Param        amount        query  number  false  "Amount"
// @Param        type          query  string  false  "Transaction Type"
// @Param        description   query  string  false  "Description"
// @Param        date          query  string  false  "Transaction Date"
// @Security     BearerAuth
// @Success      200 {object} pb.ListTransactionsResponse "List of transactions"
// @Failure      500 {string} string "Error while fetching transactions"
// @Router       /transaction/list [get]
func (h *Handler) ListTransactions(ctx *gin.Context) {
	req := &pb.ListTransactionsRequest{
		AccountId:   ctx.Query("account_id"),
		CategoryId:  ctx.Query("category_id"),
		Type:        ctx.Query("type"),
		Description: ctx.Query("description"),
		Date:        ctx.Query("date"),
	}

	// Convert the amount parameter to float32 if provided
	if amount := ctx.Query("amount"); amount != "" {
		amountValue, err := strconv.ParseFloat(amount, 32)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid amount format"})
			return
		}
		req.Amount = float32(amountValue)
	}

	res, err := h.Transaction.ListTransactions(ctx, req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}
