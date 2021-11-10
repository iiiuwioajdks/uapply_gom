package models

type Interviewers struct {
	UID            int32 `gorm:"index"`
	OrganizationID int
	DepartmentID   int `gorm:"index"`
}
