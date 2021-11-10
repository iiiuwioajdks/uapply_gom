package models

type UserWxInfo struct {
	UID        int32  `gorm:"primaryKey"`
	OpenId     string `gorm:"index"`
	SessionKey string
	Role       int8 `gorm:"comment '0是用户，1是面试官'"`
}
