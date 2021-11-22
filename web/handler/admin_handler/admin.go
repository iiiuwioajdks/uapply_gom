package admin_handler

import (
	"context"
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"uapply_go/web/forms"
	"uapply_go/web/global"
	"uapply_go/web/global/errInfo"
	"uapply_go/web/models"
	"uapply_go/web/models/response"
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
		// 判断 MySql 错误码类型，当错误码为 1062 时说明触发了 account 的唯一索引错误
		if mysqlErr, ok := result.Error.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				return errInfo.ErrDepExist
			}
		}
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
	if req.DepartmentName != "" {
		dep.DepartmentName = req.DepartmentName
	}
	if req.Account != "" {
		dep.Account = req.Account
	}
	if req.Password != "" {
		dep.Password = req.Password
	}
	// 部门 role 肯定是0
	dep.Role = 0

	//更新 Department 数据
	result := db.Model(&models.Department{}).Omit("organization_id", "department_id").Where("department_id = ? and organization_id = ?", dep.DepartmentID, dep.OrganizationID).Limit(1).Updates(&dep)
	// rowsAffected 等于 0 说明参数有误
	if result.RowsAffected == 0 {
		return sql.ErrNoRows
	}
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetDepDetail 获取部门详细信息
func GetDepDetail(depid string) (*forms.AdminReq, error) {
	db := global.DB

	// 返回值声明
	depInfo := new(forms.AdminReq)
	result := db.Table("department").Where("department_id = ? and role = 0", depid).Find(&depInfo)

	if result.Error != nil {
		return nil, result.Error
	}
	// 没有查询到结果
	if result.RowsAffected == 0 {
		return nil, sql.ErrNoRows
	}
	return depInfo, nil
}

func GetDepRoughDetail(depid int) (*forms.AdminReq, error) {
	db := global.DB

	// 返回变量声明
	depInfo := new(forms.AdminReq)
	result := db.Table("department").Omit("account", "password").Where("department_id = ?", depid).Find(&depInfo)

	if result.Error != nil {
		return nil, result.Error
	}
	// 没有查到数据
	if result.RowsAffected == 0 {
		return nil, sql.ErrNoRows
	}
	return depInfo, nil
}

func GetInterviewee(uid string, depid int, orgid int) (*models.UserInfo, error) {
	db := global.DB

	var userInfo *models.UserInfo

	// 判断用户是否报名本部门
	var count int64
	if db.Model(models.UserRegister{}).Where("organization_id = ? and department_id = ? and uid = ?", orgid, depid, uid).Count(&count); count == 0 {
		return nil, errInfo.ErrUserNotFind
	}

	// 在用户表中查找用户信息
	result := db.Model(models.UserInfo{}).Where("uid = ?", uid).First(&userInfo)
	if result.RowsAffected == 0 {
		return nil, sql.ErrNoRows
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return userInfo, nil
}

func GetAllInterviewees(orgid, depid int) ([]response.IntervieweeRsp, error) {
	db := global.DB

	var interviewees []response.IntervieweeRsp
	// 从 user_info 和 user_register 两张表里查数据
	result := db.Table("user_register as ur").Joins("inner join user_info as ui on ur.uid = ui.uid").Where("ur.organization_id = ? and ur.department_id = ?", orgid, depid).Scan(&interviewees)

	// 没有已报名的用户
	if result.RowsAffected == 0 {
		return nil, sql.ErrNoRows
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return interviewees, nil
}

// SetTime 设置报名时间
func SetTime(did int, t *forms.Time) error {
	db := global.DB
	dep := models.Department{
		StartTime: t.Start,
		EndTime:   t.End,
	}
	res := db.Model(&models.Department{}).Where("department_id=?", did).Updates(&dep)
	return res.Error
}

func GetUserInfo(depid int, orgid int) (sum int64, males int64, females int64, err error) {
	db := global.DB
	var users []*models.UserRegister
	// 本部门的报名人数,并保存uid
	result := db.Model(&models.UserRegister{}).Select("uid").Where("organization_id = ? and department_id = ?", orgid, depid).Find(&users).Count(&sum)
	if result.Error != nil {
		return 0, 0, 0, result.Error
	}
	if result.RowsAffected == 0 {
		return 0, 0, 0, sql.ErrNoRows
	}

	var flag int64
	// 报名人数中的男生数
	for _, user := range users {
		// 如果是男生，flag 变1
		flag = 0
		result := db.Model(&models.UserInfo{}).Where("uid = ? and sex = ?", user.UID, 1).Count(&flag)
		if flag == 1 {
			males++
		}
		if result.Error != nil {
			return 0, 0, 0, result.Error
		}
	}
	// 报名人数中的女生数
	females = sum - males

	return sum, males, females, nil
}
