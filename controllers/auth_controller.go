package controllers

import (
	"net/http"

	"github.com/donajivt/go-auth-service/models"
	"github.com/donajivt/go-auth-service/services"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service services.AuthService
}

func NewAuthController(s services.AuthService) *AuthController {
	return &AuthController{service: s}
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req models.RegistrationRequestDto
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ResponseDto{IsSuccess: false, Message: err.Error()})
		return
	}
	if msg, err := c.service.Register(req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ResponseDto{IsSuccess: false, Message: msg})
	} else {
		ctx.JSON(http.StatusOK, models.ResponseDto{IsSuccess: true})
	}
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req models.LoginRequestDto
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ResponseDto{IsSuccess: false, Message: err.Error()})
		return
	}
	resp, err := c.service.Login(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ResponseDto{IsSuccess: false, Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, models.ResponseDto{IsSuccess: true, Result: resp})
}

func (c *AuthController) AssignRole(ctx *gin.Context) {
	var req models.RegistrationRequestDto
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ResponseDto{IsSuccess: false, Message: err.Error()})
		return
	}
	if err := c.service.AssignRole(req.Email, req.Role); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ResponseDto{IsSuccess: false, Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, models.ResponseDto{IsSuccess: true})
}
