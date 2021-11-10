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

// CreateDep 创建部门
func CreateDep(req *forms.AdminReq) error {
	// 拿取db
	db := global.DB
	// 给model添加参数值
	var dep models.Department
	bindDepModel(req, &dep)
	// 存进数据库
	result := db.Save(&dep)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateDep 更新部门信息
// 返回值第一个参数为 RowsAffected
func UpdateDep(req *forms.AdminReq) (rowsAffected int64, err error) {
	db := global.DB
	// 给 Model 绑定参数
	var dep models.Department
	// todo:不用单独拿出来，拿出来的应该是有共同特征的语句，不然会造成结构的混乱
	bindDepModel(req, &dep)
	//更新 Department 数据
	// todo:Model里面的参数需要改成对应的类型，条件更新不要直接写实体，写成 &models.Department{}这种,然后后面跟where,不然可能会出现多余条件,确保一切在自己的掌控中是一种好习惯
	result := db.Model(&dep).Updates(&dep)
	// todo:api里面的东西拿来这里
	if result.Error != nil {
		err = result.Error
	}
	rowsAffected = result.RowsAffected
	return
}

// bindDepModel 给 Department 绑定参数
func bindDepModel(req *forms.AdminReq, dep *models.Department) {
	dep.DepartmentID = uint(req.DepartmentID)
	dep.OrganizationID = uint(req.OrganizationID)
	dep.DepartmentName = req.DepartmentName
	dep.Account = req.Account
	dep.Password = req.Password
	// 部门 role 肯定是0
	dep.Role = 0
}
