package api

import (
	"log"

	"github.com/dilshodforever/5-oyimtixon/api/handler"
	"github.com/dilshodforever/5-oyimtixon/api/middleware"
	_ "github.com/dilshodforever/5-oyimtixon/docs"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Auth Service
// @version 1.0
// @description Dilshod
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewGin(h *handler.Handler) *gin.Engine {
	ca, err := casbin.NewEnforcer("config/model.conf", "config/policy.csv")
	if err != nil {
		panic(err)
	}

	err = ca.LoadPolicy()
	if err != nil {
		log.Fatal("casbin error load policy: ", err)
		panic(err)
	}

	r := gin.Default()

	r.Use(middleware.NewAuth(ca))
	url := ginSwagger.URL("swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler, url))

	a := r.Group("/auth")
	{
		a.POST("/register", h.Register)
		a.POST("/login", h.Login)
		a.POST("/forgot-password", h.ForgotPassword)
		a.PUT("/reset-password", h.ResetPassword)
	}

	u := r.Group("/user")
	{
		u.GET("/profile", h.GetProfile)
		u.PUT("/profile", h.UpdateProfile)
		u.PUT("/change-password", h.ChangePassword)
	}

	return r
}
