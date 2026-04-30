package main

import (
	"os"

	api "github.com/gauravsahay007/split-wise-clone/api/handler"
	"github.com/gauravsahay007/split-wise-clone/auth"
	"github.com/gauravsahay007/split-wise-clone/business"
	_ "github.com/gauravsahay007/split-wise-clone/config"
	_ "github.com/gauravsahay007/split-wise-clone/docs"
	"github.com/gauravsahay007/split-wise-clone/infra"
	"github.com/gauravsahay007/split-wise-clone/middleware"
	"github.com/gauravsahay007/split-wise-clone/repository"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	auth.LoadGoogleAuthEnv(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"))
	auth.LoadGithubAuthEnv(os.Getenv("GITHUB_CLIENT_ID"), os.Getenv("GITHUB_CLIENT_SECRET"))

	db := infra.InitDB()

	repo := &repository.Repo{DB: db}

	service := &business.Service{Repo: repo}

	h := &api.Handler{Service: service}

	//Initialise gin router
	r := gin.Default()

	r.GET("/auth/:provider", h.OAuthHandler)
	r.GET("/auth/:provider/callback", h.GenerateTokenFromGoogle)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/auth/local/signup", h.UserHandler)
	r.POST("/auth/local/login", h.LoginHandler)

	authorized := r.Group("/api")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.POST("/groups", h.CreateGroupHandler)
		authorized.POST("/groups/:id/members", h.AddMemberHandler)
		authorized.POST("/expenses", h.ExpenseHandler)
		authorized.GET("/groups/:id/balances", h.BalancesHandler)
		authorized.GET("/user/summary", h.UserSummaryHandler)
		authorized.GET("/user-details", h.UserDetailsHandler)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		panic("failed to start server: " + err.Error())
	}
}
