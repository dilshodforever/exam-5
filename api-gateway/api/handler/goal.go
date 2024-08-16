package handler

import (
	"strconv"

	"github.com/dilshodforever/5-oyimtixon/api/middleware"
	pb "github.com/dilshodforever/5-oyimtixon/genprotos/goals"
	"github.com/gin-gonic/gin"
)

// CreateGoal handles creating a new goal
// @Summary      Create Goal
// @Description  Create a new goal
// @Tags         Goal
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        goal body pb.CreateGoalRequest true "Goal details"
// @Success      200 {object} pb.GoalResponse "Goal created successfully"
// @Failure      400 {string} string "Invalid input"
// @Failure      500 {string} string "Error while creating goal"
// @Router       /goal/create [post]
func (h *Handler) CreateGoal(ctx *gin.Context) {
	var req pb.CreateGoalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	id:=middleware.GetUserId(ctx)
	req.UserId=id
	res, err := h.Goal.CreateGoal(ctx, &req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}

// GetGoalByid handles retrieving a goal by ID
// @Summary      Get Goal by ID
// @Description  Retrieve details of a goal by ID
// @Tags         Goal
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Goal ID"
// @Success      200 {object} pb.GetGoalResponse "Goal details"
// @Failure      400 {string} string "Missing or invalid ID"
// @Failure      500 {string} string "Error while fetching goal"
// @Router       /goal/get/{id} [get]
func (h *Handler) GetGoalByid(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(400, gin.H{"error": "Missing or invalid ID"})
		return
	}

	req := &pb.GetGoalByidRequest{Id: id}

	res, err := h.Goal.GetGoalByid(ctx, req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}

// UpdateGoal handles updating a goal
// @Summary      Update Goal
// @Description  Update details of a goal
// @Tags         Goal
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        goal body pb.UpdateGoalRequest true "Updated goal details"
// @Success      200 {object} pb.GoalResponse "Goal updated successfully"
// @Failure      400 {string} string "Invalid input"
// @Failure      500 {string} string "Error while updating goal"
// @Router       /goal/update [put]
func (h *Handler) UpdateGoal(ctx *gin.Context) {
	var req pb.UpdateGoalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	res, err := h.Goal.UpdateGoal(ctx, &req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}

// DeleteGoal handles deleting a goal
// @Summary      Delete Goal
// @Description  Delete a goal by ID
// @Tags         Goal
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Goal ID"
// @Success      200 {object} pb.GoalResponse "Goal deleted successfully"
// @Failure      400 {string} string "Missing or invalid ID"
// @Failure      500 {string} string "Error while deleting goal"
// @Router       /goal/delete/{id} [delete]
func (h *Handler) DeleteGoal(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		ctx.JSON(400, gin.H{"error": "Missing or invalid ID"})
		return
	}

	req := &pb.DeleteGoalRequest{Id: id}

	res, err := h.Goal.DeleteGoal(ctx, req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}

// ListGoals handles listing all goals
// @Summary      List Goals
// @Description  Get a list of all goals based on the provided query parameters
// @Tags         Goal
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user_id        query     string  false  "User ID"  example("user123")
// @Param        name           query     string  false  "Name of the goal"  example("Save for vacation")
// @Param        target_amount  query     number   false  "Target amount"  example(5000.00)
// @Param        current_amount query     number   false  "Current amount"  example(1500.00)
// @Param        deadline       query     string  false  "Deadline (YYYY-MM-DD)"  example("2024-12-31")
// @Param        status         query     string  false  "Status of the goal"  example("In Progress")
// @Success      200            {object}  pb.ListGoalsResponse "List of goals"
// @Failure      500            {string}  string                  "Error while fetching goals"
// @Router       /goal/list [get]
func (h *Handler) ListGoals(ctx *gin.Context) {
	req := &pb.ListGoalsRequest{
		UserId:   ctx.Query("user_id"),
		Name:     ctx.Query("name"),
		Status:   ctx.Query("status"),
		Deadline: ctx.Query("deadline"),
	}

	// Convert target_amount and current_amount from string to float32
	if targetAmountStr := ctx.Query("target_amount"); targetAmountStr != "" {
		if targetAmount, err := strconv.ParseFloat(targetAmountStr, 32); err == nil {
			req.TargetAmount = float32(targetAmount)
		}
	}
	if currentAmountStr := ctx.Query("current_amount"); currentAmountStr != "" {
		if currentAmount, err := strconv.ParseFloat(currentAmountStr, 32); err == nil {
			req.CurrentAmount = float32(currentAmount)
		}
	}

	res, err := h.Goal.ListGoals(ctx, req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}
