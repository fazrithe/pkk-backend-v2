package controllers

import (
	"net/http"
	"strconv"

	"github.com/fazrithe/pkk-backend-v2/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InstController struct {
	DB *gorm.DB
}

func NewInstController(DB *gorm.DB) InstController {
	return InstController{DB}
}

func (inc *InstController) FindInst(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var inst []models.Institution
	result := inc.DB.Limit(intLimit).Offset(offset).Find(&inst)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(inst), "data": inst})
}
