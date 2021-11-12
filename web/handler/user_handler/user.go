package user_handler

import (
	"context"
	"database/sql"
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

func Login(code string) (token string, err error) {
	url := "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"

	// 合成url, 这里的appId和secret是在微信公众平台上获取的
	url = fmt.Sprintf(url,
		global.Conf.WXInfo.Appid,
		global.Conf.WXInfo.Secret,
		code)
	client := &http.Client{}

	request, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return "", err
	}

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var ws forms.WxSession
	if err := json.Unmarshal(body, &ws); err != nil {
		return "", err
	}
	// todo: delete
	if ws.OpenID == "" {
		return "", errInfo.ErrWXCode
	}
	log.Println("ws1:", ws)

	j := middleware.NewJWT()
	token, err = j.CreateWXToken(jwt.WXClaims{
		Role:       0,
		Openid:     ws.OpenID,
		SessionKey: ws.SessionKey,
	})
	if err != nil {
		return "", err
	}

	// 判断数据库是否存在，不存在则添加
	var userLogin models.UserWxInfo
	if result := global.DB.Where(models.UserWxInfo{OpenId: ws.OpenID}).First(&userLogin); result.RowsAffected == 1 {
		return
	}
	userLogin.Role = 0
	userLogin.SessionKey = ws.SessionKey
	userLogin.OpenId = ws.OpenID
	global.DB.Save(&userLogin)

	return
}

// GetUID 根据 OpenId 获取 UID
func GetUID(openId string) (int32, error) {
	db := global.DB

	// UserWxInfo 声明
	userWxInfo := new(models.UserWxInfo)
	result := db.Table("user_wx_info").Select("uid").Where("open_id = ?", openId).First(userWxInfo)
	if result.Error != nil {
		return 0, result.Error
	}
	// 没有查到用户信息
	if result.RowsAffected == 0 {
		return 0, sql.ErrNoRows
	}
	return userWxInfo.UID, nil
}

// SaveResume 保存用户简历
func SaveResume(req *forms.UserInfoReq) error {
	db := global.DB

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
