package repo

import (
	"gin-ddd-example/internal/app/model"
	"gin-ddd-example/pkg/db"
)

// 定义 user 仓储接口
type UserRepo interface {
	Save(user model.User) error
}

type UserRepoImpl struct {
	db *db.Database
}

func NewUserRepo(db *db.Database) *UserRepoImpl {
	return &UserRepoImpl{db}
}

// ----------- (●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●) ------------

// 持久化操作

func (repo *UserRepoImpl) Save(user model.User) error {
	result := repo.db.DB.Create(&user)
	return result.Error
}
