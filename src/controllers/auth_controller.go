package controllers

import (
	"basic-go/src/dto"
	"basic-go/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IAuthController interface {
	Signup(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type AuthController struct {
	service services.IAuthService
}

func NewAuthController(service services.IAuthService) *AuthController {
	return &AuthController{service: service}
}

func (c *AuthController) Signup(ctx *gin.Context) {
	var input dto.SignupInput

	// バリデーション
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.service.Signup(input.Email, input.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (c *AuthController) Login(ctx *gin.Context) {
	var input dto.LoginInput

	// バリデーション
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := c.service.Login(input.Email, input.Password)
	if err != nil {
		if err.Error() == "User not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
