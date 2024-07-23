package repo

import (
	"gin-ddd-example/internal/app/model"
	"gin-ddd-example/pkg/db"
)

// 定义仓储接口
type BlogRepo interface {
}

type BlogRepoImpl struct {
	*db.Database
}

func NewBlogRepo(db *db.Database) *BlogRepoImpl {
	return &BlogRepoImpl{Database: db}
}

func (repo *BlogRepoImpl) FindById(id int) (*model.Blog, error) {
	var blog model.Blog
	repo.DB.First(&blog, id)
	return &blog, nil
}

func (repo *BlogRepoImpl) Create(blog *model.Blog) (*model.Blog, error) {
	res := repo.DB.Create(blog)
	return blog, res.Error
}
