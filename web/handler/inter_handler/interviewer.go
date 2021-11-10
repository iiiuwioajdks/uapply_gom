package inter_handler

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
	"uapply_go/web/models/jwt"
)

func Login(code string) (token string, err error) {
	url := "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"

	// 合成url, 这里的appId和secret是在微信公众平台上获取的
	url = fmt.Sprintf(url,
		global.Conf.InterWxInfo.Appid,
		global.Conf.InterWxInfo.Secret,
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
		Role:       1,
		Openid:     ws.OpenID,
		SessionKey: ws.SessionKey,
	})
	if err != nil {
		return "", err
	}

	return
}

// Check 检查是否是部长（面试官，是就设数据库role为1，否则设为0）
func Check() error {
	return nil
}
