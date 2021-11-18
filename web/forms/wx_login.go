package forms

type WxSession struct {
	SessionKey string `json:"session_key"`
	ExpireIn   int64  `json:"expires_in"`
	UnionID    string `json:"unionid"`
}
