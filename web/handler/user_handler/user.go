package user_handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
	"uapply_go/web/forms"
	"uapply_go/web/global"
	"uapply_go/web/global/errInfo"
	"uapply_go/web/middleware"
	"uapply_go/web/models"
	"uapply_go/web/models/jwt"
)

const (
	Week = time.Hour * 24 * 7
)

func Login(code string) (token string, uid int32, err error) {
	url := "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"

	// 合成url, 这里的appId和secret是在微信公众平台上获取的
	url = fmt.Sprintf(url,
		global.Conf.WXInfo.Appid,
		global.Conf.WXInfo.Secret,
		code)
	client := &http.Client{}

	request, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return "", 0, err
	}

	response, err := client.Do(request)
	if err != nil {
		return "", 0, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", 0, err
	}

	var ws forms.WxSession
	if err := json.Unmarshal(body, &ws); err != nil {
		return "", 0, err
	}
	// todo: delete
	if ws.OpenID == "" {
		return "", 0, errInfo.ErrWXCode
	}
	log.Println("ws1:", ws)

	if err != nil {
		return "", 0, err
	}
	j := middleware.NewJWT()
	claim := jwt.WXClaims{ // UID自增。没有则创建时添加UID
		Role:       0,
		Openid:     ws.OpenID,
		SessionKey: ws.SessionKey,
	}
	// 判断数据库是否存在，不存在则添加
	var userLogin models.UserWxInfo
	if result := global.DB.Where(models.UserWxInfo{OpenId: ws.OpenID}).First(&userLogin); result.RowsAffected == 0 {
		userLogin.Role = 0
		userLogin.SessionKey = ws.SessionKey
		userLogin.OpenId = ws.OpenID
		global.DB.Save(&userLogin)
	}
	claim.UID = userLogin.UID
	uid = userLogin.UID
	token, err = j.CreateWXToken(claim)
	return
}

// SaveResume 保存用户简历
func SaveResume(req *forms.UserInfoReq) error {
	db := global.DB
	var count int64
	db.Model(&models.UserInfo{}).Where("uid=?", req.UID).Count(&count)
	if count != 0 {
		return errInfo.ErrResumeExist
	}

	// 绑定 model 参数
	userInfo := &models.UserInfo{
		UID:     req.UID,
		Name:    req.Name,
		StuNum:  req.StuNum,
		Address: req.Address,
		Major:   req.Major,
		Phone:   req.Phone,
		Email:   req.Email,
		Sex:     req.Sex,
		Intro:   req.Intro,
	}
	result := db.Save(userInfo)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Register 用户报名
func Register(regInfo *forms.UserRegisterInfo) error {
	db := global.DB

	// 绑定 model 参数
	reg := &models.UserRegister{
		UID:            regInfo.UID,
		OrganizationID: regInfo.OrganizationID,
		DepartmentID:   regInfo.DepartmentID,
	}
	// 保存数据库
	// 判断简历是否存在
	var count int64
	if db.Model(&models.UserInfo{}).Where(reg.UID).Count(&count); count == 0 {
		return errInfo.ErrResumeNotExist
	}

	// 判断部门是否存在
	if db.Model(models.Department{}).Where("department_id = ? and organization_id = ?", regInfo.DepartmentID, regInfo.OrganizationID).Count(&count); count == 0 {
		return errInfo.ErrInvalidParam
	}

	// 不可重复报名某一组织
	mutex := global.Rs.NewMutex(fmt.Sprintf("reg_%d", reg.UID))
	if err := mutex.Lock(); err != nil {
		return err
	}
	if db.Model(models.UserRegister{}).Where("uid=? and organization_id=?", regInfo.UID, regInfo.OrganizationID).Count(&count); count != 0 {
		return errInfo.ErrReRegister
	}
	if result := db.Create(&reg); result.Error != nil {
		return result.Error
	}
	if ok, err := mutex.Unlock(); !ok || err != nil {
		return errInfo.ErrSystem
	}
	return nil
}

// UpdateResume 更新简历
func UpdateResume(req *forms.UserResumeInfo) error {
	db := global.DB

	result := db.Model(&models.UserInfo{}).Omit("uid").Where("uid = ?", req.UID).Updates(&req)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetResume 获取简历
func GetResume(uid int32) (*forms.UserInfoReq, error) {
	db := global.DB

	resume := new(forms.UserInfoReq)
	// 查询数据
	result := db.Model(&models.UserInfo{}).Where("uid = ?", uid).First(resume)
	if result.Error != nil {
		return nil, result.Error
	}
	return resume, nil
}

// ClearText 一键清空自我介绍
func ClearText(req *forms.UserInfoReq) error {
	db := global.DB

	resume := &models.UserInfo{
		Intro: "",
	}

	// 数据库中更改intro
	result := db.Select("intro").Where("uid = ?", req.UID).Save(&resume)
	if result.Error != nil {
		return result.Error
	}

	return nil

}

// SaveTmpResume 临时保存简历
func SaveTmpResume(tmp *forms.UserResumeInfo) error {
	// 首先判断他有没有提交简历了
	rdb := global.Rdb
	exist, err := CheckIfResume(tmp.UID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if exist {
		return errInfo.ErrResumeExist
	}

	tmpByte, err := json.Marshal(tmp)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// 原子性处理
	res := rdb.SetEX(context.Background(), strconv.Itoa(int(tmp.UID)), string(tmpByte[:]), Week)
	return res.Err()
}

// CheckIfResume 检查简历是否保存
func CheckIfResume(uid int32) (bool, error) {
	db := global.DB
	res := db.Model(models.UserInfo{}).First(&models.UserInfo{}, uid)
	if res.Error != nil && res.Error.Error() != "record not found" {
		return false, res.Error
	}
	return res.RowsAffected == 1, nil
}

// GetTmpResume 获得草稿箱简历
func GetTmpResume(uid int32) (*forms.UserResumeInfo, error) {
	rdb := global.Rdb
	res := rdb.Get(context.Background(), strconv.Itoa(int(uid)))
	if res.Err() != nil {
		if res.Err() == redis.Nil {
			return nil, redis.Nil
		}
		return nil, res.Err()
	}

	var u forms.UserResumeInfo
	err := json.Unmarshal([]byte(res.Val()), &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func JudgeTime(req *forms.UserRegisterInfo) error {
	db := global.DB
	var dep models.Department
	// 先查部门
	db.First(&dep, req.DepartmentID)

	time := time.Now().Unix()

	// 如果你没设置报名开始时间，我就默认你啥时候都能报名
	if dep.StartTime == 0 {
		return nil
	}
	// 否则我就要判断你这个时间合不合法
	if time < dep.StartTime || time > dep.EndTime {
		return errInfo.ErrCNotReg
	}

	return nil
}
