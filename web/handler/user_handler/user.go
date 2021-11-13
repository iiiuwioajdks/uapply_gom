package user_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"uapply_go/web/forms"
	"uapply_go/web/global"
	"uapply_go/web/global/errInfo"
	"uapply_go/web/middleware"
	"uapply_go/web/models"
	"uapply_go/web/models/jwt"
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
	claim := jwt.WXClaims{
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
	if result := db.Model(models.UserInfo{}).Select("uid").Where("uid = ?", regInfo.UID).First(reg); result.RowsAffected == 0 {
		fmt.Println(result.Error)
		return result.Error
	}
	// 判断部门是否存在
	if result := db.Model(models.Department{}).Select("department_id").Where("department_id = ?", regInfo.DepartmentID).First(reg); result.RowsAffected == 0 {
		return result.Error
	}
	result := db.Create(reg)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateResume(req *forms.UserResumeInfo) error {
	db := global.DB

	result := db.Model(&models.UserInfo{}).Omit("uid").Where("uid = ?", req.UID).Updates(&req)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
