package inter_handler

import (
	"context"
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
