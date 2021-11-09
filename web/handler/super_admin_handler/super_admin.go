package super_admin_handler

import (
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
