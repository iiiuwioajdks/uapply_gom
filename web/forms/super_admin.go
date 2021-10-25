package forms

type CreateSAdmin struct {
	OrganizationName string `json:"organization_name" binding:"required"`
	Account          string `json:"account" binding:"required"`
	Password         string `json:"password" binding:"required"`
}
