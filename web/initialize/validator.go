package initialize

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	v1 "uapply_go/web/validator"
)

func InitValidators() {
	v := binding.Validator.Engine().(*validator.Validate)
	v.RegisterValidation("mobile", v1.ValidateMobile)
	v.RegisterValidation("email", v1.ValidateEmail)
}
