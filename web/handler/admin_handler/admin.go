package admin_handler

import (
	"context"
	"uapply_go/web/forms"
	"uapply_go/web/global"
	"uapply_go/web/global/errInfo"
	"uapply_go/web/models"
)

func Login(ctx context.Context, loginInfo *forms.Login) (*models.Department, error) {
	db := global.DB
	var admin models.Department
	result := db.Where(&models.Department{Account: loginInfo.Account, Password: loginInfo.Password}).First(&admin)
	if result.RowsAffected == 0 {
		return nil, errInfo.ErrUserNotFind
	}
	return &admin, nil
}

func CreateDep(req *forms.AdminReq) error {
	// 拿取db
	db := global.DB
	// 给model添加参数值
	var dep models.Department
	dep.DepartmentName = req.DepartmentName
	dep.Account = req.Account
	dep.Password = req.Password
	dep.Role = 0
	// 存进数据库
	result := db.Save(&dep)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
