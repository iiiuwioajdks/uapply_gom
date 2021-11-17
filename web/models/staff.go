package models

type StaffInfo struct {
	UID            int32  `json:"uid" gorm:"index"`
	UserName       string `json:"user_name"`
	Phone          string `json:"phone"`
	Year           string `json:"year"` // 入校年份，例如 2020，2021
	OrganizationID int32  `json:"organization_id" gorm:"index"`
	DepartmentID   int32  `json:"department_id" gorm:"index"`
	Role           int    `json:"role" gorm:"default:0;type:int comment '1表示部长，0表示部员'"` // 角色，部长或者部员，1表示部长，0表示部员
}
