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
