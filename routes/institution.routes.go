package routes

import (
	"github.com/fazrithe/pkk-backend-v2/controllers"
	"github.com/fazrithe/pkk-backend-v2/middleware"
	"github.com/gin-gonic/gin"
)

type InstitutionRouteController struct {
	institutionController controllers.InstitutionController
}

func NewRouteInstitutionController(institutionController controllers.InstitutionController) InstitutionRouteController {
	return InstitutionRouteController{institutionController}
}

func (ic *InstitutionRouteController) InstitutionRoute(rg *gin.RouterGroup) {
	router := rg.Group("institutions")
	router.Use(middleware.DeserializeUser())
	router.POST("", ic.institutionController.CreateInstitution)
	router.GET("", ic.institutionController.FindInstitutions)
	router.GET("/:institutionId", ic.institutionController.FindInstitutionById)
	router.PUT("/:institutionId", ic.institutionController.UpdateInstitution)
	router.DELETE("/:institutionId", ic.institutionController.DeleteController)
	router.GET("/select", ic.institutionController.SelectInstitution)
}
