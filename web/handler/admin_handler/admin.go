package admin_handler

import (
	"context"
	"uapply_go/web/forms"
	"uapply_go/web/global"
	"uapply_go/web/global/errInfo"
	"uapply_go/web/models"
)

func Login(ctx context.Context, loginInfo *forms.Login) (error, *models.Department) {
	db := global.DB
	var admin models.Department
	result := db.Where(&models.Department{Account: loginInfo.Account, Password: loginInfo.Password}).First(&admin)
	if result.RowsAffected == 0 {
		return errInfo.ErrUserNotFind, nil
	}
	return nil, &admin
}
