package routes

import (
	"net/http"

	"github.com/fazrithe/pkk-backend-v2/controllers"
	"github.com/fazrithe/pkk-backend-v2/middleware"
	"github.com/gin-gonic/gin"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewRouteUserController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (uc *UserRouteController) UserRoute(rg *gin.RouterGroup) {

	router := rg.Group("users")
	router.GET("/me", middleware.DeserializeUser(), uc.userController.GetMe)
	router.GET("", middleware.DeserializeUser(), uc.userController.FindUser)
	router.POST("", middleware.DeserializeUser(), uc.userController.CreateUser)
	router.GET("/:userId", middleware.DeserializeUser(), uc.userController.FindUserById)
	router.PUT("/:userId", middleware.DeserializeUser(), uc.userController.UpdateUser)
	router.DELETE("/:userId", middleware.DeserializeUser(), uc.userController.DeleteUser)
	router.POST("/upload", middleware.DeserializeUser(), uc.userController.UploadPhoto)
	router.StaticFS("/images", http.Dir("public"))
}
