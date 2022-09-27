package main

import (
	"belajarbwa/auth"
	"belajarbwa/campaign"
	"belajarbwa/handler"
	"belajarbwa/helper"
	"belajarbwa/user"
	"fmt"
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
	// user repository, service, handler
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := handler.NewUserHandler(userService, authService)
	// campaign repository, service, handler
	campaignRepository := campaign.NewRepository(db)
	campaigns, err := campaignRepository.FindByUserID(1)
	fmt.Println(len(campaigns))
	for _, campaign := range campaigns {
		fmt.Println(campaign.Name)
		if len(campaign.CampaignImages) > 0 {
			fmt.Println(campaign.CampaignImages[0].FileName)
		}
	}
	// router
	router := gin.Default()
	// version 1
	api := router.Group("/api/v1")
	// group users
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)

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
