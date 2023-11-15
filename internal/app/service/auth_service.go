package service

import "gin-ddd-example/internal/app/repo"

type AuthService interface {
	Signup(req *SignupDto) error
}

type AuthServiceImpl struct {
	userRepo repo.UserRepo
}

func NewAuthService(userRepo repo.UserRepo) *AuthServiceImpl {
	return &AuthServiceImpl{userRepo}
}

type SignupDto struct {
	Name     string `form:"name" json:"name" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required"`
}

func (s *AuthServiceImpl) Signup(req *SignupDto) error {
	return nil
}
