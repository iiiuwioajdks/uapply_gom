package models

type UserWxInfo struct {
	UID        int32  `gorm:"primaryKey"`
	OpenId     string `gorm:"index"`
	SessionKey string
	Role       int `json:"role" gorm:"default:0;type:int comment '1表示面试官，0表示用户'"` // 角色，部长或者部员，1表示面试官，0表示用户
}
