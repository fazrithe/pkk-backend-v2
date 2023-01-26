package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/fazrithe/pkk-backend-v2/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InstitutionController struct {
	DB *gorm.DB
}

func NewInstitutionController(DB *gorm.DB) InstitutionController {
	return InstitutionController{DB}
}

func (ic *InstitutionController) CreateInstitution(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var payload *models.CreateInstitutionRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	now := time.Now()
	newInstitution := models.Institution{
		Name:      payload.Name,
		Address:   payload.Address,
		Image:     payload.Image,
		User:      currentUser.ID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := ic.DB.Create(&newInstitution)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Institution with that title already exists"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newInstitution})
}
