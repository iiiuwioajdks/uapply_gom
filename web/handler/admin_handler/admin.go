package admin_handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"sync"
	"time"
	"uapply_go/web/forms"
	"uapply_go/web/global"
	"uapply_go/web/global/errInfo"
	"uapply_go/web/models"
	jwt2 "uapply_go/web/models/jwt"
	"uapply_go/web/models/response"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111" // 引入sms
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
	rdb := global.Rdb
	// 判断此子是否为该部门的
	key := fmt.Sprintf("inter_%d", uid.UID)
	redisRes := rdb.SetNX(context.Background(), key, 1, time.Second*2)
	if redisRes.Val() == false {
		return errInfo.ErrConcurrent
	}
	tx := db.Begin()
	res := tx.Model(models.StaffInfo{}).Where(&models.StaffInfo{DepartmentID: int32(id.DepartmentID), OrganizationID: int32(id.OrganizationID), UID: int32(uid.UID)}).
		Update("role", 1)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errInfo.ErrInvalidParam
	}
	// 将 user_wx_info 的role设置为1
	res = tx.Model(&models.UserWxInfo{}).Where(`uid=?`, uid.UID).Update("role", 1)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}
	// 添加到interviewers
	res = tx.Create(&models.Interviewers{
		UID:            int32(uid.UID),
		OrganizationID: id.OrganizationID,
		DepartmentID:   id.DepartmentID,
	})
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}
	tx.Commit()
	return nil
}

func GetPhones(depid int, orgid int, uids *forms.MultiUIDForm) (*response.PhoneInfo, error) {
	db := global.DB
	// 剔除非法的uid
	var phoneInfo response.PhoneInfo
	var ids []int
	var phones []string
	sqlRaw := "select uid from user_register where uid in (?) and organization_id = ? and department_id = ?"
	if res := db.Raw(sqlRaw, uids.UID, orgid, depid).Find(&ids); res.Error != nil {
		return nil, res.Error
	}
	sqlRaw = "select phone from user_info where uid in (?)"
	if res := db.Raw(sqlRaw, uids.UID).Find(&phones); res.Error != nil {
		return nil, res.Error
	}
	if len(phones) != len(ids) {
		return nil, errInfo.ErrInvalidParam
	}
	phoneInfo.UID = ids
	phoneInfo.Phone = phones
	return &phoneInfo, nil
}

func Out(form *forms.MultiUIDForm, orgid, depid int) error {
	db := global.DB
	result := db.Model(&models.UserRegister{}).Where("organization_id = ? and department_id = ? and uid in ?", orgid, depid, form.UID).Delete(&models.UserRegister{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

var enrollLock sync.Mutex

func Enroll(form *forms.MultiUIDForm, orgid, depid int) error {
	db := global.DB
	// 可能会有部分 uid 已经存在在 user_enroll 表中了，需要在后面添加数据的时候将这些排除在外，防止重复添加
	var uidExists []int
	// 去重后的 uid 切片声明
	var uidDistinct []int
	enrollLock.Lock()
	defer enrollLock.Unlock()
	// 查找存在的 uid
	result := db.Table("user_enroll").Select("uid").Where("organization_id = ? and department_id = ? and uid in ?", orgid, depid, form.UID).Find(&uidExists)
	// 当查到数据时，进行去重
	if result.RowsAffected != 0 {
		uidDistinct = make([]int, 0, len(form.UID))
		tmap := make(map[int]bool)
		for _, uid := range uidExists {
			tmap[uid] = true
		}
		for _, uid := range form.UID {
			// map 里没有对应的 key 说明没有重复
			if _, ok := tmap[uid]; !ok {
				uidDistinct = append(uidDistinct, uid)
			}
		}
	} else {
		// 没查到就不用去重
		uidDistinct = form.UID
	}

	// 查找经过去重后的用户信息
	var userEnroll []models.UserEnroll
	sqlRawSelect := "select ur.uid,ui.name as user_name,organization_id,department_id from user_register as ur inner join user_info ui on ur.uid = ui.uid where organization_id = ? and department_id = ? and ur.uid in ?"
	db.Raw(sqlRawSelect, orgid, depid, uidDistinct).Scan(&userEnroll)
	// 先更新 user_register 里用户的状态，然后再将对应用户信息添加到 user_enroll 表中
	sqlRawUpdate := "update user_register set first_status=3,second_status=3,final_status=1 where organization_id = ? and department_id = ? and uid in ?"
	tx := db.Begin()
	result = tx.Exec(sqlRawUpdate, orgid, depid, uidDistinct)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if len(userEnroll) > 0 {
		result = tx.Table("user_enroll").Save(&userEnroll)
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}
	tx.Commit()
	return nil
}

func SendSMS(Type int, pi *response.PhoneInfo) error {
	credential := common.NewCredential(
		global.Conf.SMSInfo.SecretId,
		global.Conf.SMSInfo.SecretKey,
	)
	cpf := profile.NewClientProfile()

	cpf.HttpProfile.ReqMethod = "POST"

	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"

	cpf.SignMethod = "HmacSHA1"

	client, _ := sms.NewClient(credential, "ap-guangzhou", cpf)

	request := sms.NewSendSmsRequest()

	/* 短信应用ID: 短信SdkAppId在 [短信控制台] 添加应用后生成的实际SdkAppId，示例如1400006666 */
	request.SmsSdkAppId = common.StringPtr(global.Conf.SMSInfo.SdkAppId)
	/* 短信签名内容: 使用 UTF-8 编码，必须填写已审核通过的签名，签名信息可登录 [短信控制台] 查看 */
	request.SignName = common.StringPtr("xxx")
	/* 用户的 session 内容: 可以携带用户侧 ID 等上下文信息，server 会原样返回 */
	//request.SessionContext = common.StringPtr("xxx")
	/* 模板参数: 若无模板参数，则设置为空*/
	var model string
	switch Type {
	case 1: // 通过第一轮
		request.TemplateParamSet = common.StringPtrs([]string{""})
		model = global.Conf.SMSInfo.Model1
	case 2: // 通过第二轮
		request.TemplateParamSet = common.StringPtrs([]string{""})
		model = global.Conf.SMSInfo.Model2
	case 3: // 录取
		request.TemplateParamSet = common.StringPtrs([]string{""})
		model = global.Conf.SMSInfo.Model3
	case 4: // 淘汰
		request.TemplateParamSet = common.StringPtrs([]string{""})
		model = global.Conf.SMSInfo.Model4
	default:
		return errInfo.ErrInvalidParam
	}
	/* 模板 ID: 必须填写已审核通过的模板 ID。模板ID可登录 [短信控制台] 查看 */
	request.TemplateId = common.StringPtr(model)
	/* 下发手机号码，采用 E.164 标准，+[国家或地区码][手机号]
	 * 示例如：+8613711112222， 其中前面有一个+号 ，86为国家码，13711112222为手机号，最多不要超过200个手机号*/
	for i := 0; i < len(pi.Phone); i++ {
		pi.Phone[i] = "+86" + pi.Phone[i]
	}
	request.PhoneNumberSet = common.StringPtrs(pi.Phone)

	// 通过client对象调用想要访问的接口，需要传入请求对象
	response, err := client.SendSms(request)
	// 处理异常
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return err
	}
	if err != nil {
		return err
	}
	b, _ := json.Marshal(response.Response)
	zap.S().Info(b)
	return nil
}

func GetInterviewed(num string, orgid int, depid int) ([]*models.UserInfo, error) {
	db := global.DB

	var intervieweds []*models.UserInfo
	if num == "1" {
		sqlRaw := "SELECT ui.`uid`, name FROM user_register ur JOIN user_info ui ON ui.`uid` = ur.`uid` WHERE organization_id = ? AND department_id = ? AND (first_status = 2 OR first_status = 3)"
		result := db.Raw(sqlRaw, orgid, depid).Find(&intervieweds)
		if result.Error != nil {
			return nil, result.Error
		}
	} else if num == "2" {
		sqlRaw := "SELECT ui.`uid`, name FROM user_register ur JOIN user_info ui ON ui.`uid` = ur.`uid` WHERE organization_id = ? AND department_id = ? AND first_status = 3 AND (second_status = 2 OR second_status = 3)"
		result := db.Raw(sqlRaw, orgid, depid).Find(&intervieweds)
		if result.Error != nil {
			return nil, result.Error
		}

	} else {
		return nil, errInfo.ErrInvalidParam
	}

	return intervieweds, nil
}

// GetUserEnroll 部门获取自己的通过部员
func GetUserEnroll(orgid int, depid int) ([]*models.UserInfo, error) {
	db := global.DB

	var enrolls []*models.UserInfo
	sqlRaw := "SELECT ui.`uid`, name FROM user_register ur JOIN user_info ui ON ui.`uid` = ur.`uid` WHERE organization_id = ? AND department_id = ? AND final_status = 1"
	result := db.Raw(sqlRaw, orgid, depid).Find(&enrolls)
	if result.Error != nil {
		return nil, result.Error
	}
	return enrolls, nil

}
func Pass(num string, orgid int, depid int, uids forms.MultiUIDForm) error {
	db := global.DB

	if num == "1" {
		sqlRaw := "update user_register set first_status = 3 where deleted_at IS NULL and department_id = ? and organization_id = ? and uid IN (?)"
		result := db.Exec(sqlRaw, depid, orgid, uids.UID)
		if result.Error != nil {
			return result.Error
		}

	} else if num == "2" {
		sqlRaw := "update user_register set second_status = 3 where deleted_at IS NULL and department_id = ? and organization_id = ? and uid IN (?) and first_status = 3"
		result := db.Exec(sqlRaw, depid, orgid, uids.UID)
		if result.Error != nil {
			return result.Error
		}
	} else {
		return errInfo.ErrInvalidParam
	}

	return nil
}
