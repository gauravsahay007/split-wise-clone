package main

import (
	"os"

	"github.com/gauravsahay007/split-wise-clone/api"
	"github.com/gauravsahay007/split-wise-clone/business"
	"github.com/gauravsahay007/split-wise-clone/infra"
	"github.com/gauravsahay007/split-wise-clone/repository"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	db := infra.InitDB()

	repo := &repository.Repo{DB: db}

	service := &business.Service{Repo: repo}

	h := &api.Handler{Service: service}

	//Initialise gin router
	r := gin.Default()

	v1 := r.Group("/api")
	{
		v1.POST("/users", h.UserHandler)
		v1.POST("/expenses", h.ExpenseHandler)
		v1.POST("/groups", h.CreateGroupHandler)
		v1.POST("/groups/:id/members", h.AddMemberHandler)
		v1.GET("/groups/:id/balances", h.BalancesHandler)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
