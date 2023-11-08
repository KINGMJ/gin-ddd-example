package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	fmt.Println("登录请求")
	c.JSON(http.StatusOK, gin.H{"message": "key", "status": 200})
}

type SignupDto struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" validate:"required"`
}

func Signup(c *gin.Context) {
	var signupDto SignupDto
	if err := c.ShouldBind(&signupDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprint(err)})
		return
	}

	// err =

	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}
