package validate

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// 注册自定义的校验
func RegisterCustomValidate() {
	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validate.RegisterValidation("phone", validatePhone)
	}
}

func validatePhone(fl validator.FieldLevel) bool {
	phone, ok := fl.Field().Interface().(string)
	if ok {
		phoneRegex := regexp.MustCompile(`^[\d-]+$`)
		return phoneRegex.MatchString(phone)
	}
	return false
}
