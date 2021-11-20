package super_admin_handler

import (
	"database/sql"
	"uapply_go/web/forms"
	"uapply_go/web/global"
	"uapply_go/web/models"
)

// Create 创建超级管理员
func Create(csa *forms.CreateSAdmin) error {
	db := global.DB
	// 存 Organization 的信息
	orgModel := models.Organization{
		OrganizationName: csa.OrganizationName,
	}
	result := db.Save(&orgModel)
	if result.Error != nil {
		return result.Error
	}

	// 存 Organization 的信息 到department表，作用是存储账号密码，方便登录查询
	csaModel := models.Department{
		OrganizationID: orgModel.OrganizationID,
		// name "-admin"
		DepartmentName: orgModel.OrganizationName + "-admin",
		Account:        csa.Account,
		Password:       csa.Password,
		Role:           1,
	}
	result = db.Save(&csaModel)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Update 更新超级管理员信息
func Update(req *forms.UpdateSAdmin) error {
	db := global.DB
	// 更新 Organization 的信息
	if req.OrganizationName != "" {
		var orgModel models.Organization
		orgModel.OrganizationName = req.OrganizationName
		result := db.Model(&models.Organization{}).Omit("organization_id").Where("organization_id = ?", req.OrganizationID).
			Updates(&orgModel)
		// rowsAffected 等于 0 说明参数有误
		if result.RowsAffected == 0 {
			return sql.ErrNoRows
		}
		if result.Error != nil {
			return result.Error
		}
	}

	if req.Account != "" || req.Password != "" {
		var depModel models.Department
		depModel.Account = req.Account
		depModel.Password = req.Password
		result := db.Model(&models.Department{}).Omit("organization_id", "department_id").
			Where("organization_id=? and role=1", req.OrganizationID).
			Updates(depModel)
		if result.RowsAffected == 0 {
			return sql.ErrNoRows
		}
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func DeleteDepartment(depid string, orgid int) error {
	//获取数据库
	db := global.DB
	//数据库操作
	result := db.Where("department_id = ?", depid).Where("organization_id = ? and role = 0", orgid).Delete(&models.Department{})

	// 可能是组织id错了，要删别的组织的部门，直接退回说参数错误
	// 也可能是部门id错了，那也一样是参数错误
	// 如果删除的是组织，那更不行，要考虑全部，包括爬虫请求在内
	if result.RowsAffected == 0 {
		return sql.ErrNoRows
	}

	//数据库的查询错误
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func ShowConcreteDepartInfo(depid string, orgid int) (*forms.AdminReq, error) {
	// 获取DB
	db := global.DB

	// 返回变量的声明
	departInfo := new(forms.AdminReq)

	// 数据库操作
	// Find 返回的是影响 不是错误（err）
	result := db.Table("department").Where("department_id = ?", depid).Where("organization_id = ? and role = 0", orgid).Find(departInfo)
	if result.RowsAffected == 0 {
		// 可能是这个部门不在这个组织
		// 部门id错误
		// 不要把超级管理员的信息展示
		return nil, sql.ErrNoRows
	}
	return departInfo, nil
}

func GetOrganizationInfo(orgid string) (*models.Organization, error) {
	// 获取数据库
	db := global.DB

	// 从数据库中获取组织信息
	var org models.Organization
	result := db.Select("organization_name", "created_at").Where("organization_id = ?", orgid).First(&org)

	// 数据库的查询错误
	if result.Error != nil {
		return nil, result.Error
	}

	return &org, nil
}

func GetOrgDepartments(orgid int) ([]*models.Department, error) {
	// 获取数据库
	db := global.DB

	// 从数据库中获取部门名称
	var deps []*models.Department
	result := db.Select("department_name", "department_id").Where("organization_id = ? and role = 0", orgid).Find(&deps)

	// 数据库的查询错误
	if result.Error != nil {
		return nil, result.Error
	}

	return deps, nil
}

func SetTime(orgid int, t *forms.Time) error {
	db := global.DB
	dep := models.Department{
		StartTime: t.Start,
		EndTime:   t.End,
	}
	res := db.Model(&models.Department{}).Where("organization_id=?", orgid).Updates(&dep)
	return res.Error
}
