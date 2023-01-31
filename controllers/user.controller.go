package controllers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fazrithe/pkk-backend-v2/models"
	"github.com/fazrithe/pkk-backend-v2/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(DB *gorm.DB) UserController {
	return UserController{DB}
}

func (uc *UserController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	userResponse := &models.UserResponse{
		ID:        currentUser.ID,
		Name:      currentUser.Name,
		Email:     currentUser.Email,
		Photo:     currentUser.Photo,
		Role:      currentUser.Role,
		Provider:  currentUser.Provider,
		CreatedAt: currentUser.CreatedAt,
		UpdatedAt: currentUser.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": userResponse}})
}

func (uc *UserController) FindUser(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var users []models.User
	result := uc.DB.Limit(intLimit).Offset(offset).Find(&users)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Fail", "message": result.Error})
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(users), "data": users})
}

func (uc *UserController) CreateUser(ctx *gin.Context) {
	var payload *models.SignUpInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if payload.Password != payload.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Password do not match"})
		return
	}

	now := time.Now()

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	newUser := models.User{
		Name:      payload.Name,
		Email:     strings.ToLower(payload.Email),
		Password:  hashedPassword,
		Role:      payload.Role,
		Verified:  true,
		Photo:     payload.Photo,
		Provider:  "local",
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := uc.DB.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key") {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User with the email already"})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Something bad happened"})
		return
	}

	userResponse := &models.UserResponse{
		ID:        newUser.ID,
		Name:      newUser.Name,
		Email:     newUser.Email,
		Photo:     newUser.Photo,
		Role:      newUser.Role,
		Provider:  newUser.Provider,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": userResponse}})

}

func (uc *UserController) FindUserById(ctx *gin.Context) {
	userId := ctx.Param("userId")

	var user models.User
	result := uc.DB.First(&user, "id = ?", userId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No user with that name"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": user})
}

func (uc *UserController) UpdateUser(ctx *gin.Context) {
	userId := ctx.Param("userId")

	var payload *models.SignUpInput
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var updateUser models.User
	result := uc.DB.First(&updateUser, "id = ?", userId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No user with that name"})
		return
	}
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}
	now := time.Now()
	userToUpdate := models.User{
		Name:      payload.Name,
		Email:     payload.Email,
		Role:      payload.Role,
		Photo:     payload.Photo,
		Password:  hashedPassword,
		CreatedAt: updateUser.CreatedAt,
		UpdatedAt: now,
	}

	uc.DB.Model(&updateUser).Updates(userToUpdate)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": updateUser})
}

func (uc *UserController) DeleteUser(ctx *gin.Context) {
	userId := ctx.Param("userId")

	result := uc.DB.Delete(&models.User{}, "id = ?", userId)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No user with that name"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (uc *UserController) UploadPhoto(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("photo")
	if err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}

	fileExt := filepath.Ext(header.Filename)
	originalFileName := strings.TrimSuffix(filepath.Base(header.Filename), filepath.Ext(header.Filename))
	now := time.Now()
	filename := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", now.Unix()) + fileExt
	filePath := "http://localhost:8080/api/users/images/users/" + filename

	out, err := os.Create("public/users/" + filename)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}
	ctx.JSON(http.StatusOK, gin.H{"filepath": filePath})
}
