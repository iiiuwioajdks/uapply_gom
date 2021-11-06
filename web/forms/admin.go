package forms

type Login struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AdminReq 为了有较好的复用性，不在此处做binding要求，因此在代码中需要判定参数是否符合条件
type AdminReq struct {
	DepartmentID   int    `json:"department_id"`
	OrganizationID int    `json:"organization_id"`
	DepartmentName string `json:"department_name"`
	Account        string `json:"account"`
	Password       string `json:"password"`
}
