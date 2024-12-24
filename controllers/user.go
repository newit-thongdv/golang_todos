package controllers

import (
	"golang-todos/auth"
	"golang-todos/database"
	"golang-todos/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

/**
* Login
**/
func Login(context *gin.Context) {
	var request TokenRequest
	var user models.User

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	record := database.Instance.Where("email = ?", request.Email).First(&user)
	if record.Error != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}

	credentialError := user.CheckPassword(request.Password)
	if credentialError != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		context.Abort()
		return
	}

	accessTokenString, err := auth.GenerateJWT(user.Email, auth.ExpiredTime)
	refreshTokenString, errRefresh := auth.GenerateJWT(user.Email, auth.ExpiredRefreshTime)

	if err != nil || errRefresh != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	user.RefreshToken = refreshTokenString
	if err := database.Instance.Save(&user).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update refresh token"})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"access_token":  accessTokenString,
		"expired_Time":  auth.ExpiredTime,
		"refresh_token": refreshTokenString})
}

/**
* Handle refresh token
**/
func HandleRefresh(context *gin.Context) {
	var user models.User
	tokens, err := auth.RefreshToken(user.RefreshToken)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new access token"})
		return
	}

	user.RefreshToken = tokens.RefreshToken
	if err := database.Instance.Save(&user).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update refresh token"})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"access_token":  tokens.AccessToken,
		"expired_Time":  auth.ExpiredTime,
		"refresh_token": tokens.RefreshToken})
}

/**
* Register user
**/
func RegisterUser(context *gin.Context) {
	var user models.User
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if err := user.HashPassword(user.Password); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	record := database.Instance.Create(&user)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}
	context.JSON(http.StatusCreated, gin.H{"userId": user.ID, "email": user.Email})
}
