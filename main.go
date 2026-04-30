package main

import (
	"os"
	"time"

	api "github.com/gauravsahay007/split-wise-clone/api/handler"
	"github.com/gauravsahay007/split-wise-clone/auth"
	"github.com/gauravsahay007/split-wise-clone/business"
	_ "github.com/gauravsahay007/split-wise-clone/config"
	_ "github.com/gauravsahay007/split-wise-clone/docs"
	"github.com/gauravsahay007/split-wise-clone/infra"
	"github.com/gauravsahay007/split-wise-clone/middleware"
	"github.com/gauravsahay007/split-wise-clone/repository"
	"github.com/gin-contrib/cors"
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

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", os.Getenv("FRONTEND_URL")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/auth/:provider", h.OAuthHandler)
	r.GET("/auth/:provider/callback", h.GenerateTokenFromGoogle)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/auth/local/signup", h.UserHandler)
	r.POST("/auth/local/login", h.LoginHandler)

	authorized := r.Group("/api")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.POST("/create-group", h.CreateGroupHandler)
		authorized.POST("/groups/:id/members", h.AddMemberHandler)
		authorized.POST("/add-expenses", h.ExpenseHandler)
		authorized.POST("/groups/:id/expenses", h.GetGroupExpenses)
		authorized.GET("/groups/:id/balances", h.BalancesHandler)
		authorized.GET("/user-summary", h.UserSummaryHandler)
		authorized.GET("/user-details", h.UserDetailsHandler)
		authorized.GET("/groups", h.GetUserGroupsHandler)
		authorized.POST("/add-friends", h.HandleAddFriend)
		authorized.GET("/get-friends", h.GetFriendsList)
		authorized.GET("/search-friends", h.SearchFriendsInAGroup)
		authorized.GET("/group/members", h.GetGroupMembers)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		panic("failed to start server: " + err.Error())
	}
}
