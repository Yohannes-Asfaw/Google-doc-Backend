package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/khallihub/godoc/dto"
	"github.com/khallihub/godoc/service"
)

type SignupController interface {
	Signup(ctx *gin.Context) string
}

type signupController struct {
	signupService service.SignupService
}

func NewSignupController(signupService service.SignupService) SignupController {
	return &signupController{
		signupService: signupService,
	}
}

func (controller *signupController) Signup(ctx *gin.Context) string {
	var signUPInfo dto.SignUp

	err := ctx.ShouldBind(&signUPInfo)
	if err != nil {
		return ""
	}

	err = controller.signupService.Signup(signUPInfo.FullName, signUPInfo.Email, signUPInfo.Password)
	if err != nil {
		return "Username already exists"
	}

	// If needed, return a success message or handle the response accordingly
	return "User created successfully"
}
