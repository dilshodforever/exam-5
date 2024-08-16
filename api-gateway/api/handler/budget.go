package handler

import (
	"strconv"

	"github.com/dilshodforever/5-oyimtixon/api/middleware"
	pb "github.com/dilshodforever/5-oyimtixon/genprotos/budgets"
	"github.com/dilshodforever/5-oyimtixon/kafkasender"
	"github.com/gin-gonic/gin"
)

// CreateBudget handles creating a new budget
// @Summary      Create Budget
// @Description  Create a new budget
// @Tags         Budget
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        budget body pb.CreateBudgetRequest true "Budget details"
// @Success      200 {object} pb.BudgetResponse "Budget created successfully"
// @Failure      400 {string} string "Invalid input"
// @Failure      500 {string} string "Error while creating budget"
// @Router       /budget/create [post]
func (h *Handler) CreateBudget(ctx *gin.Context) {
	var req pb.CreateBudgetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	id := middleware.GetUserId(ctx)
	req.UserId=id
	res, err := h.Budget.CreateBudget(ctx, &req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, res)
}

// GetBudgetByid handles retrieving a budget by ID
// @Summary      Get Budget by ID
// @Description  Get details of a budget by ID
// @Tags         Budget
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Budget ID"
// @Success      200 {object} pb.GetBudgetByidResponse "Budget details"
// @Failure      400 {string} string "Missing or invalid ID"
// @Failure      500 {string} string "Error while fetching budget"
// @Router       /budget/get/{id} [get]
func (h *Handler) GetBudgetByid(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		ctx.JSON(400, gin.H{"error": "Missing or invalid ID"})
		return
	}

	req := &pb.GetBudgetByidRequest{Id: id}

	res, err := h.Budget.GetBudgetbyid(ctx, req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}

// UpdateBudget handles updating a budget
// @Summary      Update Budget
// @Description  Update details of a budget
// @Tags         Budget
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        budget body pb.UpdateBudgetRequest true "Updated budget details"
// @Success      200 {object} pb.BudgetResponse "Budget updated successfully"
// @Failure      400 {string} string "Invalid input"
// @Failure      500 {string} string "Error while updating budget"
// @Router       /budget/update [put]
func (h *Handler) UpdateBudget(ctx *gin.Context) {
	var req pb.UpdateBudgetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	res, err := kafkasender.UpdateBudget(h.Kafka, &req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res.Message)
}

// DeleteBudget handles deleting a budget
// @Summary      Delete Budget
// @Description  Delete a budget by ID
// @Tags         Budget
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Budget ID"
// @Success      200 {object} pb.BudgetResponse "Budget deleted successfully"
// @Failure      400 {string} string "Missing or invalid ID"
// @Failure      500 {string} string "Error while deleting budget"
// @Router       /budget/delete/{id} [delete]
func (h *Handler) DeleteBudget(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(400, gin.H{"error": "Missing or invalid ID"})
		return
	}

	req := &pb.DeleteBudgetRequest{Id: id}

	res, err := h.Budget.DeleteBudget(ctx, req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}

// ListBudgets handles listing all budgets
// @Summary      List Budgets
// @Description  Get a list of all budgets based on the provided query parameters
// @Tags         Budget
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user_id    query     string  false  "User ID"  example("user123")
// @Param        category_id query    string  false  "Category ID"  example("category456")
// @Param        amount     query     number   false  "Amount"  example(1500.75)
// @Param        period     query     string  false  "Period"  example("monthly")
// @Param        start_date query     string  false  "Start date (YYYY-MM-DD)"  example("2024-01-01")
// @Param        end_date   query     string  false  "End date (YYYY-MM-DD)"  example("2024-12-31")
// @Success      200       {object}  pb.ListBudgetsResponse "List of budgets"
// @Failure      500       {string}  string                 "Error while fetching budgets"
// @Router       /budget/list [get]
func (h *Handler) ListBudgets(ctx *gin.Context) {
	req := &pb.ListBudgetsRequest{
		UserId:     ctx.Query("user_id"),
		CategoryId: ctx.Query("category_id"),
		Period:     ctx.Query("period"),
		StartDate:  ctx.Query("start_date"),
		EndDate:    ctx.Query("end_date"),
	}

	if amountStr := ctx.Query("amount"); amountStr != "" {
		if amount, err := strconv.ParseFloat(amountStr, 32); err == nil {
			req.Amount = float32(amount)
		}
	}

	res, err := h.Budget.ListBudgets(ctx, req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}
