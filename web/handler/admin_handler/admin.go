package admin_handler

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"time"
	"uapply_go/web/forms"
	"uapply_go/web/global"
	"uapply_go/web/global/errInfo"
	"uapply_go/web/models"
	jwt2 "uapply_go/web/models/jwt"
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

func GetAllInterviewees(depid, orgid int) ([]response.IntervieweeRsp, error) {
	db := global.DB

	var interviewees []response.IntervieweeRsp
	// 从 user_info 和 user_register 两张表里查数据
	sqlRaw := "SELECT * FROM user_register as ur inner join user_info as ui on ur.uid = ui.uid where ur.organization_id = ? and ur.department_id = ?"
	result := db.Raw(sqlRaw, orgid, depid).Scan(&interviewees)

	// 没有已报名的用户
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

func GetUserInfo(depid int, orgid int) (rsp *response.FMInfo, err error) {
	db := global.DB
	var users []*models.UserRegister
	rsp = &response.FMInfo{}
	// 本部门的报名人数,并保存uid
	result := db.Model(&models.UserRegister{}).Select("uid").Where("organization_id = ? and department_id = ?", orgid, depid).Find(&users).Count(&rsp.Sum)
	if result.Error != nil {
		return rsp, result.Error
	}

	// 报名人数中的男生数
	sqlRaw := "SELECT COUNT(sex) FROM user_info as ui inner join user_register as ur on ur.uid = ui.uid where ur.organization_id = ? and ur.department_id = ? and ui.sex=1"
	db.Raw(sqlRaw, orgid, depid).Scan(&rsp.Male)
	rsp.Female = rsp.Sum - rsp.Male
	return rsp, nil
}

func AddInterviewers(id *jwt2.Claims, uid *forms.Interviewer) error {
	db := global.DB
	// 判断此子是否为该部门的
	rdb := global.Rdb
	key := fmt.Sprintf("inter_%d", uid.UID)
	redisRes := rdb.SetNX(context.Background(), key, 1, time.Second*2)
	if redisRes.Val() == false {
		return errInfo.ErrConcurrent
	}

	res := db.Model(models.StaffInfo{}).Where(&models.StaffInfo{DepartmentID: int32(id.DepartmentID), OrganizationID: int32(id.OrganizationID), UID: int32(uid.UID)}).
		Update("role", 1)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errInfo.ErrInvalidParam
	}
	// 将 user_wx_info 的role设置为1
	res = db.Model(&models.UserWxInfo{}).Where(`uid=?`, uid.UID).Update("role", 1)
	if res.Error != nil {
		return res.Error
	}
	// 添加到interviewers
	res = db.Create(&models.Interviewers{
		UID:            int32(uid.UID),
		OrganizationID: id.OrganizationID,
		DepartmentID:   id.DepartmentID,
	})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func Out(form *forms.MultiUIDForm, orgid, depid int) error {
	// 开启事务
	tx := global.DB.Begin()
	result := tx.Table("user_register").Where("organization_id = ? and department_id = ? and uid in ?", orgid, depid, form.UID).Delete(&models.UserRegister{})
	// RowsAffected 和 uid 切片长度不一致说明有部分 uid 不正确
	if result.RowsAffected != int64(len(form.UID)) {
		// 回滚
		tx.Rollback()
		return errInfo.ErrInvalidUIDS
	}
	if result.Error != nil {
		// 回滚
		tx.Rollback()
		return result.Error
	}
	tx.Commit()
	return nil
}
