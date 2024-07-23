package repo_test

import (
	"fmt"
	"gin-ddd-example/internal/app/model"
	"gin-ddd-example/internal/app/repo"
	"gin-ddd-example/pkg/utils"
	"testing"
)

func TestAutoMigrate(t *testing.T) {
	database.AutoMigrate(&model.Blog{})
}

func TestFindByID(t *testing.T) {
	repo := &repo.BlogRepoImpl{database}
	res, err := repo.FindById(7)
	if err != nil {
		fmt.Println(err)
	}
	utils.PrettyJson(res)
}
