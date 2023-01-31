package routes

import (
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
	router.Use(middleware.DeserializeUser())
	router.GET("/me", uc.userController.GetMe)
	router.GET("", uc.userController.FindUser)
	router.POST("", uc.userController.CreateUser)
	router.GET("/:userId", uc.userController.FindUserById)
	router.PUT("/:userId", uc.userController.UpdateUser)
	router.DELETE("/:userId", uc.userController.DeleteUser)
}
