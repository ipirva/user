package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/vmwarecloudadvocacy/user/internal/auth"
	"github.com/vmwarecloudadvocacy/user/internal/db"
	"github.com/vmwarecloudadvocacy/user/internal/service"
	"github.com/vmwarecloudadvocacy/user/pkg/logger"
)

const (
	dbName         = "acmefit"
	collectionName = "users"
)

// This handles initiation of "gin" router. It also defines routes to various APIs
// Env variable USER_IP and USER_PORT should be used to set IP and PORT.
// For example: export USER_PORT=8086 will start the server on local IP at 0.0.0.0:8086
func handleRequest() {

	// Init Router
	router := gin.New()

	nonAuthGroup := router.Group("/")
	{
		nonAuthGroup.POST("/register", service.RegisterUser)
		nonAuthGroup.POST("/login", service.LoginUser)
		nonAuthGroup.POST("/refresh-token", service.RefreshAccessToken)
		nonAuthGroup.POST("/verify-token", service.VerifyAuthToken)
	}

	authGroup := router.Group("/")
	// Added
	authGroup.Use(auth.AuthMiddleware())
	{
		authGroup.GET("/users", service.GetUsers)
		authGroup.GET("/users/:id", service.GetUser)
		authGroup.DELETE("/users/:id", service.DeleteUser)
		authGroup.POST("/logout", service.LogoutUser)
	}

	//flag.Parse()

	// Set default values if ENV variables are not set
	port := db.GetEnv("USERS_PORT", "8081")
	ip := db.GetEnv("USERS_HOST", "0.0.0.0")

	ipPort := ip + ":" + port

	logger.Logger.Infof("Starting user service at %s on %s", ip, port)

	router.Run(ipPort)

}

func main() {

	//create your file with desired read/write permissions
	f, err := os.OpenFile("log.info", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Could not open file ", err)
	} else {
		logger.InitLogger(f)
	}

	dbsession := db.ConnectDB(dbName, collectionName, logger.Logger)
	logger.Logger.Infof("Successfully connected to database %s", dbName)

	redisClient := db.ConnectRedisDB(logger.Logger)
	logger.Logger.Infof("Successfully connected to redis database NAME")

	handleRequest()

	db.CloseDB(dbsession, logger.Logger)

	defer f.Close()
	defer redisClient.Close()

}