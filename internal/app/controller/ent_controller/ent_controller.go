package ent_controller

import (
	"fmt"
	"gin-ddd-example/internal/app/model"
	"gin-ddd-example/internal/app/service/ent_service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "list ents", "status": 200})
}

func Create(c *gin.Context) {
	var addEntReq model.AddEntReq
	if err := c.ShouldBind(&addEntReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprint(err)})
		return
	}
	err := ent_service.CreateEnt(&addEntReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprint(err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ent created", "status": 200})
}

func Update(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "update ent", "status": 200})
}

func Delete(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "delete ent", "status": 200})
}
