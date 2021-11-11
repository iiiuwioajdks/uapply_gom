package admin_handler

import (
	"context"
	"database/sql"
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

// CreateDep 创建部门
func CreateDep(req *forms.AdminReq) error {
	// 拿取db
	db := global.DB
	// 给model添加参数值
	var dep models.Department
	dep.OrganizationID = uint(req.OrganizationID)
	dep.DepartmentName = req.DepartmentName
	dep.Account = req.Account
	dep.Password = req.Password
	// 部门 role 肯定是0
	dep.Role = 0

	// 存进数据库
	result := db.Save(&dep)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateDep 更新部门信息
func UpdateDep(req *forms.AdminReq) error {
	db := global.DB
	// 给 Model 绑定参数
	var dep models.Department
	dep.DepartmentID = uint(req.DepartmentID)
	dep.OrganizationID = uint(req.OrganizationID)
	dep.DepartmentName = req.DepartmentName
	dep.Account = req.Account
	dep.Password = req.Password
	// 部门 role 肯定是0
	dep.Role = 0

	//更新 Department 数据
	result := db.Model(&models.Department{}).Where("department_id = ?", dep.DepartmentID).Updates(&dep)
	// rowsAffected 等于 0 说明参数有误
	if result.RowsAffected == 0 {
		return sql.ErrNoRows
	}
	if result.Error != nil {
		return result.Error
	}
	return nil
}
