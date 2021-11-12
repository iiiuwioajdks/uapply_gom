package models

type StaffInfo struct {
	UID            int32  `json:"uid" gorm:"index"`
	UserName       string `json:"user_name"`
	Phone          string `json:"phone"`
	Year           string `json:"year"` // 入校年份，例如 2020，2021
	OrganizationID int32  `json:"organization_id" gorm:"index"`
	DepartmentID   int32  `json:"department_id" gorm:"index"`
	Role           int8   `json:"role"` // 角色，部长或者部员
}
