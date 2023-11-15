package controller

import (
	"fmt"
	"gin-ddd-example/internal/app/service"
	"gin-ddd-example/pkg/errors/e"
	"gin-ddd-example/pkg/response"

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
// @Param		 page query  int false "页码"
// @Param        page_size query int false "每页查询的数量"
// @Success      200  {object} response.Response
// @Failure      400  {object} response.Response
// @Router       /v1/ents [get]
func (c *EntController) List(ctx *gin.Context) {
	response.JSON(ctx, e.Success, "加载企业列表")
}

// @Summary      创建企业
// @Description  创建企业
// @Tags         ents
// @Accept       json
// @Produce      json
// @Param        ent body service.AddEntDto true "创建企业参数"
// @Success      200  {object} response.Response
// @Failure      400  {object} response.Response
// @Router       /v1/ents [post]
func (c *EntController) Create(ctx *gin.Context) {
	var addEntDto service.AddEntDto
	if err := ctx.ShouldBind(&addEntDto); err != nil {
		response.JSON(ctx, e.ValidateErr, fmt.Sprint(err))
		return
	}
	err := c.entService.CreateEnt(&addEntDto)
	if err != nil {
		response.JSON(ctx, e.BadRequestErr, fmt.Sprint(err))
	}
	response.JSON(ctx, e.Success, "企业创建成功")
}
