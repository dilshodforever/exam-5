package handler

import (
	"github.com/dilshodforever/5-oyimtixon/api/middleware"
	pb "github.com/dilshodforever/5-oyimtixon/genprotos/categories"
	"github.com/gin-gonic/gin"
)

// CreateCategory handles creating a new category
// @Summary      Create Category
// @Description  Create a new category
// @Tags         Category
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        category body pb.CreateCategoryRequest true "Category details"
// @Success      200 {object} pb.CategoryResponse "Category created successfully"
// @Failure      400 {string} string "Invalid input"
// @Failure      500 {string} string "Error while creating category"
// @Router       /category/create [post]
func (h *Handler) CreateCategory(ctx *gin.Context) {
	var req pb.CreateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	id := middleware.GetUserId(ctx)
	req.UserId=id
	res, err := h.Category.CreateCategory(ctx, &req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}

// UpdateCategory handles updating a category
// @Summary      Update Category
// @Description  Update details of a category
// @Tags         Category
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        category body pb.UpdateCategoryRequest true "Updated category details"
// @Success      200 {object} pb.CategoryResponse "Category updated successfully"
// @Failure      400 {string} string "Invalid input"
// @Failure      500 {string} string "Error while updating category"
// @Router       /category/update [put]
func (h *Handler) UpdateCategory(ctx *gin.Context) {
	var req pb.UpdateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	res, err := h.Category.UpdateCategory(ctx, &req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}

// DeleteCategory handles deleting a category
// @Summary      Delete Category
// @Description  Delete a category by ID
// @Tags         Category
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Category ID"
// @Success      200 {object} pb.CategoryResponse "Category deleted successfully"
// @Failure      400 {string} string "Missing or invalid ID"
// @Failure      500 {string} string "Error while deleting category"
// @Router       /category/delete/{id} [delete]
func (h *Handler) DeleteCategory(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(400, gin.H{"error": "Missing or invalid ID"})
		return
	}

	req := &pb.DeleteCategoryRequest{Id: id}

	res, err := h.Category.DeleteCategory(ctx, req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}

// ListCategories handles listing all categories
// @Summary      List Categories
// @Description  Get a list of all categories based on the provided query parameters
// @Tags         Category
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user_id  query     string  false  "User ID"  example("user123")
// @Param        name     query     string  false  "Name of the category"  example("Groceries")
// @Param        type     query     string  false  "Type of the category"  example("Expense")
// @Success      200      {object}  pb.ListCategoriesResponse "List of categories"
// @Failure      500      {string}  string                  "Error while fetching categories"
// @Router       /category/list [get]
func (h *Handler) ListCategories(ctx *gin.Context) {
	req := &pb.ListCategoriesRequest{
		UserId: ctx.Query("user_id"),
		Name:   ctx.Query("name"),
		Type:   ctx.Query("type"),
	}

	res, err := h.Category.ListCategories(ctx, req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}

// GetByidCategory handles retrieving a category by ID
// @Summary      Get Category by ID
// @Description  Get details of a category by ID
// @Tags         Category
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Category ID"
// @Success      200 {object} pb.GetByidCategoriesResponse "Category details"
// @Failure      400 {string} string "Missing or invalid ID"
// @Failure      500 {string} string "Error while fetching category"
// @Router       /category/get/{id} [get]
func (h *Handler) GetByidCategory(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(400, gin.H{"error": "Missing or invalid ID"})
		return
	}

	req := &pb.GetByidCategoriesRequest{Id: id}

	res, err := h.Category.GetByidCatagory(ctx, req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}
