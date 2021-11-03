package models

import (
	"gorm.io/gorm"
	"time"
)

type DepartmentBase struct {
	DepartmentID uint `gorm:"primarykey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// Department 部门表
type Department struct {
	// 直接用id作为departmentId，这样可以索引最优化
	DepartmentBase
	OrganizationID uint   `gorm:"index;type:int;not null"`
	DepartmentName string `gorm:"type:varchar(100);not null"`
	Account        string `gorm:"unique;type:varchar(20);not null"`
	Password       string `gorm:"type:varchar(20);not null"`
	Role           int8   `gorm:"type:int;not null comment '0代表管理员，1代表超级管理员'"`
}

type OrganizationBase struct {
	OrganizationID uint `gorm:"primarykey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

// Organization 组织表
type Organization struct {
	// 直接用id作为organizationId，这样可以索引最优化
	OrganizationBase
	OrganizationName string `gorm:"unique;type:varchar(100);not null"`
	Department       []Department
}
