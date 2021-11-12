package forms

type CreateSAdmin struct {
	OrganizationName string `json:"organization_name" binding:"required"`
	Account          string `json:"account" binding:"required"`
	Password         string `json:"password" binding:"required"`
}
type UpdateSAdmin struct {
	DepartmentID     int    `json:"department_id"`
	OrganizationID   int    `json:"organization_id"`
	OrganizationName string `json:"organization_name" binding:"required"`
	Account          string `json:"account" binding:"required"`
	Password         string `json:"password" binding:"required"`
}
