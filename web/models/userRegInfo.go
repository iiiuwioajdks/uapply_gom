package models

import "gorm.io/gorm"

type UserInfo struct {
	UID     int32  `json:"uid" gorm:"primaryKey"`
	Name    string `json:"name"` // 姓名
	StuNum  string `json:"stu_num" gorm:"index"`
	Address string `json:"address"` // 楼号
	Major   string `json:"major"`   // 专业
	Phone   string `json:"phone"`   // 手机
	Email   string `json:"email"`   // 邮箱
	Sex     int8   `json:"sex"`     // 性别， 1为男，2为女
	Intro   string `json:"intro"`   // 简介
}

type UserRegister struct {
	UID            int32          `json:"uid" gorm:"index"` // 这里的uid是可以重复的，因为用户可以报名多个组织
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	OrganizationID int            `json:"organization_id"`
	DepartmentID   int            `json:"department_id" gorm:"index"`
	FirstStatus    int8           `json:"first_status" gorm:"comment '0表示新投递，1表示已安排面试，2代表已面试，3代表面试通过'"`          // 0表示新投递，1表示已安排面试，2代表已面试，3代表面试通过
	SecondStatus   int8           `json:"second_status" gorm:"comment '0表示通过第一轮但未安排面试，1表示已安排面试，2代表已面试，3代表面试通过'"` // 0表示通过第一轮但未安排面试，1表示已安排面试，2代表已面试，3代表面试通过
	FinalStatus    int8           `json:"final_status" gorm:"comment 0表示面试进行中，1表示录取，2表示淘汰"`                      // 0表示面试进行中，1表示录取，2表示淘汰
}

type UserEnroll struct {
	UID            int32 `json:"uid"`
	UserName       string
	DepartmentID   int `gorm:"index"`
	OrganizationID int `gorm:"index"`
}
