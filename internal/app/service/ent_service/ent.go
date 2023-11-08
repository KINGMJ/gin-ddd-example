package ent_service

import (
	"gin-ddd-example/internal/app/dao/ent_dao"
	"gin-ddd-example/internal/app/model"
)

func CreateEnt(addEntReq *model.AddEntReq) error {
	ent := addEntReq.ToEnt()
	return ent_dao.Save(ent)
}
