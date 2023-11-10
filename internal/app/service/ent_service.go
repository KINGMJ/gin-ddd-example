package service

import (
	"gin-ddd-example/internal/app/model"
	"gin-ddd-example/internal/app/repo"
)

// ent 服务接口
type EntService interface {
	CreateEnt(req *AddEntDto) error
	ListEnts(page int, pageSize int) []EntListDto
	// ViewEnt(entId int) EntListDto
}

// entServiceImpl 实现EntService接口
type EntServiceImpl struct {
	entRepo repo.EntRepo
}

func NewEntService(entRepo repo.EntRepo) *EntServiceImpl {
	return &EntServiceImpl{entRepo}
}

// ----------- (●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●) ------------
// 定义各种dto

// 创建企业dto
type AddEntDto struct {
	EntName      string `form:"ent_name" json:"ent_name" binding:"required"`
	EntDesc      string `form:"ent_desc" json:"ent_desc"`
	ContactName  string `form:"contact_name" json:"contact_name" binding:"required"`
	ContactEmail string `form:"contact_email" json:"contact_email" binding:"required,email"`
	ContactPhone string `form:"contact_phone" json:"contact_phone" binding:"required,phone"`
}

// dto 与 po 转换 ，这里暂时不使用领域模型的概念
func (dto AddEntDto) ToEnt() model.Ent {
	return model.Ent{
		Name: dto.EntName,
		Desc: dto.EntDesc,
	}
}

// 企业列表dto
type EntListDto struct {
	EntName      string
	EntDesc      string
	ContactName  string
	ContactEmail string
	ContactPhone string
}

// ----------- (●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●) ------------
// 业务操作
func (s *EntServiceImpl) ListEnts(page, pageSize int) []EntListDto {
	// entPo, err := s.entRepo.FindById()
	return []EntListDto{}
}

func (s *EntServiceImpl) CreateEnt(addEntDto *AddEntDto) error {
	entEntity := addEntDto.ToEnt()
	return s.entRepo.Save(entEntity)
}

// func UpdateEnt(entId int, updateEntReq *model.UpdateEntReq) error {
// 	ent, err := ent_dao.FindById(entId)
// 	if err != nil {
// 		return fmt.Errorf("ent_id 不存在")
// 	}
// 	if updateEntReq.EntName != "" {
// 		ent.Name = updateEntReq.EntName
// 	}
// 	if updateEntReq.EntDesc != "" {
// 		ent.Desc = updateEntReq.EntDesc
// 	}
// 	return ent_dao.Update(ent)
// }
