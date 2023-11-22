package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/khallihub/godoc/dto"
	"github.com/khallihub/godoc/service"
)

type LoginController interface {
	Login(ctx *gin.Context) string
}

type loginController struct {
	loginService service.LoginService
	jWtService   service.JWTService
}

func NewLoginController(loginService service.LoginService,
	jWtService service.JWTService) LoginController {
	return &loginController{
		loginService: loginService,
		jWtService:   jWtService,
	}
}

func (controller *loginController) Login(ctx *gin.Context) string {
	var credentials dto.Login
	err := ctx.ShouldBind(&credentials)
	if err != nil {
		return ""
	}
	isAuthenticated := controller.loginService.Login(credentials.Email, credentials.Password)
	if isAuthenticated {
		return controller.jWtService.GenerateToken(credentials.Email, true)
	}
	return ""
}
