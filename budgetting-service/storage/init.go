package storage

import (
	"context"

	accpb "github.com/dilshodforever/5-oyimtixon/genprotos/accaunts"
	budgpb "github.com/dilshodforever/5-oyimtixon/genprotos/budgets"
	catpb "github.com/dilshodforever/5-oyimtixon/genprotos/categories"
	goalp "github.com/dilshodforever/5-oyimtixon/genprotos/goals"
	transpb "github.com/dilshodforever/5-oyimtixon/genprotos/transactions"
)

type InitRoot interface {
	Account() AccountService
	Budget() BudgetService
	Category() CategoryService
	Goal() GoalService
	Transaction() TransactionService
}

type AccountService interface {
	CreateAccount(ctx context.Context, req *accpb.CreateAccountRequest) (*accpb.CreateAccountResponse, error)
	GetAccountByid(ctx context.Context, req *accpb.GetByIdAccauntRequest) (*accpb.GetAccountByidResponse, error)
	UpdateAccount(ctx context.Context, req *accpb.UpdateAccountRequest) (*accpb.UpdateAccountResponse, error)
	DeleteAccount(ctx context.Context, req *accpb.DeleteAccountRequest) (*accpb.UpdateAccountResponse, error)
	ListAccounts(ctx context.Context, req *accpb.ListAccountsRequest) (*accpb.ListAccountsResponse, error)
	UpdateBalance(ctx context.Context, accountID string, amount float32) error
	UpdateBalanceMinus(ctx context.Context, accountID string, amount float32) error 
	CheckBudget(ctx context.Context, userId string) (bool, error) 
}

type BudgetService interface {
	CreateBudget(ctx context.Context, req *budgpb.CreateBudgetRequest) (*budgpb.BudgetResponse, error)
	GetBudgetByid(ctx context.Context, req *budgpb.GetBudgetByidRequest) (*budgpb.GetBudgetByidResponse, error)
	UpdateBudget(ctx context.Context, req *budgpb.UpdateBudgetRequest) (*budgpb.BudgetResponse, error)
	DeleteBudget(ctx context.Context, req *budgpb.DeleteBudgetRequest) (*budgpb.BudgetResponse, error)
	ListBudgets(ctx context.Context, req *budgpb.ListBudgetsRequest) (*budgpb.ListBudgetsResponse, error)
	UpdateBudgetAmount(ctx context.Context, UserId string, amount float32)  error
}

type CategoryService interface {
	CreateCategory(ctx context.Context, req *catpb.CreateCategoryRequest) (*catpb.CategoryResponse, error)
	UpdateCategory(ctx context.Context, req *catpb.UpdateCategoryRequest) (*catpb.CategoryResponse, error)
	DeleteCategory(ctx context.Context, req *catpb.DeleteCategoryRequest) (*catpb.CategoryResponse, error)
	ListCategories(ctx context.Context, req *catpb.ListCategoriesRequest) (*catpb.ListCategoriesResponse, error)
	GetByidCategory(ctx context.Context, req *catpb.GetByidCategoriesRequest) (*catpb.GetByidCategoriesResponse, error)
}

type GoalService interface {
	CreateGoal(ctx context.Context, req *goalp.CreateGoalRequest) (*goalp.GoalResponse, error)
	GetGoalByid(ctx context.Context, req *goalp.GetGoalByidRequest) (*goalp.GetGoalResponse, error)
	UpdateGoal(ctx context.Context, req *goalp.UpdateGoalRequest) (*goalp.GoalResponse, error)
	DeleteGoal(ctx context.Context, req *goalp.DeleteGoalRequest) (*goalp.GoalResponse, error)
	ListGoals(ctx context.Context, req *goalp.ListGoalsRequest) (*goalp.ListGoalsResponse, error)
	UpdateGoulAmount(ctx context.Context, UserId string, amount float32) error
	CheckGoal(ctx context.Context, userId string) (bool, string, error)
}

type TransactionService interface {
	CreateTransaction(ctx context.Context, req *transpb.CreateTransactionRequest) (*transpb.TransactionResponse, error)
	GetTransaction(ctx context.Context, req *transpb.GetTransactionRequest) (*transpb.GetTransactionResponse, error)
	UpdateTransaction(ctx context.Context, req *transpb.UpdateTransactionRequest) (*transpb.TransactionResponse, error)
	DeleteTransaction(ctx context.Context, req *transpb.DeleteTransactionRequest) (*transpb.TransactionResponse, error)
	ListTransactions(ctx context.Context, req *transpb.ListTransactionsRequest) (*transpb.ListTransactionsResponse, error)
}
