package routes

import (
	"github.com/fazrithe/pkk-backend-v2/controllers"
	"github.com/fazrithe/pkk-backend-v2/middleware"
	"github.com/gin-gonic/gin"
)

type InstRouteController struct {
	instController controllers.InstController
}

func NewRouteInstController(instController controllers.InstController) InstRouteController {
	return InstRouteController{instController}
}

func (inc *InstRouteController) InstRoute(rg *gin.RouterGroup) {
	router := rg.Group("inst")
	router.GET("", middleware.DeserializeUser(), inc.instController.FindInst)
}
