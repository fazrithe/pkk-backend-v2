package main

import (
	"log"
	"net/http"

	"github.com/fazrithe/pkk-backend-v2/controllers"
	"github.com/fazrithe/pkk-backend-v2/initializers"
	"github.com/fazrithe/pkk-backend-v2/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	server              *gin.Engine
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	InstController      controllers.InstController
	InstRouteController routes.InstRouteController

	InstitutionController      controllers.InstitutionController
	InstitutionRouteController routes.InstitutionRouteController
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)

	AuthController = controllers.NewAuthController(initializers.DB)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	UserController = controllers.NewUserController(initializers.DB)
	UserRouteController = routes.NewRouteUserController(UserController)

	InstController = controllers.NewInstController(initializers.DB)
	InstRouteController = routes.NewRouteInstController(InstController)

	InstitutionController = controllers.NewInstitutionController(initializers.DB)
	InstitutionRouteController = routes.NewRouteInstitutionController(InstitutionController)

	server = gin.Default()
}

func main() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8080", config.ClientOrigin}
	// corsConfig.AllowAllOrigins = true
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Accept-Encoding", "X-CSRF-Token"}
	corsConfig.AllowCredentials = true
	corsConfig.AllowMethods = []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

	server.Use(cors.New(corsConfig))

	router := server.Group("api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		message := "Welcome to Golang with Gorm and Postgres"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	AuthRouteController.AuthRoute(router)
	UserRouteController.UserRoute(router)
	InstRouteController.InstRoute(router)
	InstitutionRouteController.InstitutionRoute(router)
	log.Fatal(server.Run(":" + config.ServerPort))
}
