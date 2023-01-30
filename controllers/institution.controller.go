package controllers

import (
	"net/http"
	"strconv"
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

func (ic *InstitutionController) FindInstitutions(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var institutions []models.Institution
	result := ic.DB.Limit(intLimit).Offset(offset).Find(&institutions)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(institutions), "data": institutions})
}

func (ic *InstitutionController) SelectInstitution(ctx *gin.Context) {
	var institution *models.Institution

	// var payload *models.CreateInstitutionRequest

	result := ic.DB.Find(&institution)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
		return
	}

	selectInstitution := models.InstitutionResponse{
		Name:    institution.Name,
		Address: institution.Address,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": selectInstitution})
}

func (ic *InstitutionController) FindInstitutionById(ctx *gin.Context) {
	institutionId := ctx.Param("institutionId")

	var institution models.Institution
	result := ic.DB.First(&institution, "id = ?", institutionId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No Insitution with that name exists"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": institution})

}

func (ic *InstitutionController) UpdateInstitution(ctx *gin.Context) {
	institutionId := ctx.Param("institutionId")
	currentUser := ctx.MustGet("currentUser").(models.User)

	var payload *models.UpdateInstitution
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var updateInstitution models.Institution
	result := ic.DB.First(&updateInstitution, "id = ?", institutionId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No Institution with that name exists"})
		return
	}

	now := time.Now()
	institutionToUpdate := models.Institution{
		Name:      payload.Name,
		Address:   payload.Address,
		Image:     payload.Image,
		User:      currentUser.ID,
		CreatedAt: updateInstitution.CreatedAt,
		UpdatedAt: now,
	}

	ic.DB.Model(&updateInstitution).Updates(institutionToUpdate)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": updateInstitution})
}

func (ic *InstitutionController) DeleteController(ctx *gin.Context) {
	institutionId := ctx.Param("institutionId")

	result := ic.DB.Delete(&models.Institution{}, "id = ?", institutionId)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "Fail", "message": "No Institution with that name exists"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
