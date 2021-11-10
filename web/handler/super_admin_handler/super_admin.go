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
