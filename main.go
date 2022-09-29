package main

import (
	"belajarbwa/auth"
	"belajarbwa/campaign"
	"belajarbwa/handler"
	"belajarbwa/helper"
	"belajarbwa/user"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// koneksi db
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root:@tcp(127.0.0.1:3306)/belajar_bwa?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// klo koneksi db ada error
	if err != nil {
		log.Fatal(err.Error())
	}
	// auth
	authService := auth.NewService()
	// repository
	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)
	// service
	userService := user.NewService(userRepository)
	campaignService := campaign.NewService(campaignRepository)
	// handler
	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	// router
	router := gin.Default()
	// biar images nya bisa di buka di url
	router.Static("/images", "./images")
	// version 1
	api := router.Group("/api/v1")
	// group users
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)

	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)

	router.Run()
}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// dapetin header auhorization
		authHeader := c.GetHeader("Authorization")
		// cari kata Bearer di dalam string
		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			// berentiin karna ini middleware
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		tokenString := ""
		// Bearer tokennya
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}
		// cek validasi token
		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			// berentiin karna ini middleware
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			// berentiin karna ini middleware
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		// defaut bawaan jwt nya float64 di konvert ke int
		userID := int(claim["user_id"].(float64))
		// cek user by id
		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			// berentiin karna ini middleware
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		// set session ini
		c.Set("currentUser", user)
	}
}
