package ent_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "list ents", "status": 200})
}

func Create(c *gin.Context) {

	// var signupDto SignupDto
	// if err := c.ShouldBind(&signupDto); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprint(err)})
	// 	return
	// }

	// err =

	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})

	c.JSON(http.StatusOK, gin.H{"message": "create ent", "status": 200})
}

func Update(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "update ent", "status": 200})
}

func Delete(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "delete ent", "status": 200})
}
