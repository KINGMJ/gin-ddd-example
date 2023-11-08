package ent_service

import (
	"fmt"
	"gin-ddd-example/internal/app/dao/ent_dao"
	"gin-ddd-example/internal/app/model"
)

func CreateEnt(addEntReq *model.AddEntReq) error {
	ent := addEntReq.ToEnt()
	return ent_dao.Save(ent)
}

func UpdateEnt(entId int, updateEntReq *model.UpdateEntReq) error {
	ent, err := ent_dao.FindById(entId)
	if err != nil {
		return fmt.Errorf("ent_id 不存在")
	}
	if updateEntReq.EntName != "" {
		ent.Name = updateEntReq.EntName
	}
	if updateEntReq.EntDesc != "" {
		ent.Desc = updateEntReq.EntDesc
	}
	return ent_dao.Update(ent)
}
