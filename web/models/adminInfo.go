package models

import "gorm.io/gorm"

type Department struct {
	gorm.Model
	DepartmentName string `gorm:"type:varchar(100);not null"`
	DepartmentID   int    `gorm:"index:idx_depid;unique;type:int;not null;AUTO_INCREMENT"`
	OrganizationID int    `gorm:"index:idx_orgid;unique;type:int;not null"`
	Account        string `gorm:"unique;type:varchar(20);not null"`
	Password       string `gorm:"type:varchar(20);not null"`
}

type Organization struct {
	gorm.Model
	OrganizationName string `gorm:"type:varchar(100);not null"`
	OrganizationID   int    `gorm:"index:idx_orgid;unique;type:int;not null;AUTO_INCREMENT"`
}
