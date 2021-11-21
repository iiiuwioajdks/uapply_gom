package inter_handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		global.Conf.InterWxInfo.Appid,
		global.Conf.InterWxInfo.Secret,
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
	claim.Role = userLogin.Role
	claim.UID = userLogin.UID
	uid = userLogin.UID
	token, err = j.CreateWXToken(claim)
	return
}

// Check 检查是否是部长（面试官，是就设数据库role为1，否则设为0）
func Check() error {
	return nil
}

// GetUser 查询用户的简历
func GetUser(useUid string, interUid int32) (*models.UserInfo, error) {
	db := global.DB

	// 查找用户的报名时所填报的组织和部门
	var userAuth models.UserRegister
	if result := db.Model(&models.UserRegister{}).Where("uid = ?", useUid).First(&userAuth); result.RowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	// 查找面试官所在的组织和部门
	// 中间件已经证明这个是面试官了 肯定在表里面不用判断了
	var interAuth []models.Interviewers
	if result := db.Model(&models.Interviewers{}).Where("uid = ?", interUid).Find(&interAuth); result.Error != nil {
		// 查询错误
		return nil, result.Error
	}

	// 一个人承担多个组织或部门的面试官
	// 确定这个用户是这个面试官来面试（确定组织id和部门id是否相同）
	var flag bool = false
	for _, interviewers := range interAuth {
		if interviewers.OrganizationID == userAuth.OrganizationID && interviewers.DepartmentID == userAuth.DepartmentID {
			flag = true
			break
		}
	}

	if flag {
		var useMsg models.UserInfo
		// 能来面试前面已经判断过简历表中有信息 不用再做判断
		result := db.Model(&models.UserInfo{}).Where("uid = ?", userAuth.UID).First(&useMsg)
		if result.Error != nil {
			// 查询错误
			return nil, result.Error
		}
		return &useMsg, nil
	} else {
		// 这个用户报名的时候不是这个部门或组织
		return nil, errInfo.ErrUserMatch
	}
}
