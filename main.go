package main

import (
	"os"

	api "github.com/gauravsahay007/split-wise-clone/api/handler"
	"github.com/gauravsahay007/split-wise-clone/business"
	_ "github.com/gauravsahay007/split-wise-clone/docs"
	"github.com/gauravsahay007/split-wise-clone/infra"
	"github.com/gauravsahay007/split-wise-clone/middleware"
	"github.com/gauravsahay007/split-wise-clone/repository"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	if err := godotenv.Load(); err != nil {
		return
	}

	db := infra.InitDB()

	repo := &repository.Repo{DB: db}

	service := &business.Service{Repo: repo}

	h := &api.Handler{Service: service}

	//Initialise gin router
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/api/users", h.UserHandler)
	r.POST("/api/login", h.LoginHandler)

	authorized := r.Group("/api")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.POST("/groups", h.CreateGroupHandler)
		authorized.POST("/groups/:id/members", h.AddMemberHandler)
		authorized.POST("/expenses", h.ExpenseHandler)
		authorized.GET("/groups/:id/balances", h.BalancesHandler)
		authorized.GET("/user/summary", h.UserSummaryHandler)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		panic("failed to start server: " + err.Error())
	}
}
