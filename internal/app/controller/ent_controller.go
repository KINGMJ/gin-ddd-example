package controller

import (
	"fmt"
	"gin-ddd-example/internal/app/service"

	"github.com/gin-gonic/gin"
)

// 控制器
type EntController struct {
	entService service.EntService
}

func NewEntController(entService service.EntService) *EntController {
	return &EntController{entService}
}

// @Summary      企业列表加载
// @Description  加载我可以看到的所有企业
// @Tags         ents
// @Accept       json
// @Produce      json
// @Success      200  {string}  json{"code", "message"}
// @Router       /v1/ents [get]
func (c *EntController) List(ctx *gin.Context) {
	// ents := c.entService.
	// 	c.JSON(http.StatusOK, gin.H{"message": "list ents", "status": 200})
	fmt.Println("加载企业列表")
}

// @Summary      创建企业
// @Description  创建企业
// @Tags         ents
// @Accept       json
// @Produce      json
// @Param        ent body model.AddEntReq true "create ent"
// @Success      200  {string}  json{"code", "message"}
// @Failure      500  {string}  json{"status", "message"}
// @Router       /v1/ents [post]
// func Create(c *gin.Context) {
// 	var addEntReq model.AddEntReq
// 	if err := c.ShouldBind(&addEntReq); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprint(err)})
// 		return
// 	}
// 	err := ent_service.CreateEnt(&addEntReq)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprint(err)})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "ent created", "status": 200})
// }

// func Update(c *gin.Context) {
// 	var updateEntReq model.UpdateEntReq
// 	if err := c.ShouldBind(&updateEntReq); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprint(err)})
// 		return
// 	}
// 	entId, _ := strconv.Atoi(c.Param("id"))
// 	err := ent_service.UpdateEnt(entId, &updateEntReq)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprint(err)})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "ent updated", "status": 200})
// }

// func Delete(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{"message": "delete ent", "status": 200})
// }
