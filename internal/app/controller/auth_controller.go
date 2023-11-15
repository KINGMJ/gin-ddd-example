package controller

import (
	"fmt"
	"gin-ddd-example/internal/app/service"
	"gin-ddd-example/pkg/errors/e"
	"gin-ddd-example/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService}
}

// @Summary      用户注册
// @Description  用户注册
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        signupData body service.SignupDto true "用户注册参数"
// @Success      200  {object} response.Response
// @Failure      400  {object} response.Response
// @Router       /auth/signup [post]
func (c *AuthController) Signup(ctx *gin.Context) {
	var signupDto service.SignupDto
	if err := ctx.ShouldBind(&signupDto); err != nil {
		response.JSON(ctx, e.ValidateErr, fmt.Sprint(err))
		return
	}
	err := c.authService.Signup(&signupDto)
	if err != nil {
		response.JSON(ctx, e.BadRequestErr, fmt.Sprint(err))
	}
	response.JSON(ctx, e.Success, "用户注册成功")
}
