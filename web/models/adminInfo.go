package models

import "gorm.io/gorm"

type Department struct {
	gorm.Model
	OrganizationID uint   `gorm:"index:idx_orgid;unique;type:int;not null"`
	DepartmentName string `gorm:"type:varchar(100);not null"`
	Account        string `gorm:"unique;type:varchar(20);not null"`
	Password       string `gorm:"type:varchar(20);not null"`
	Role           int8   `gorm:"type:int;not null comment '0代表管理员，1代表超级管理员'"`
}

type Organization struct {
	gorm.Model
	OrganizationName string `gorm:"type:varchar(100);not null"`
	Department       []Department
}
