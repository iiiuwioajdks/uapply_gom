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

// UIDForm 像AddExtraEnroll这种其实只要一个uid就够了
type UIDForm struct {
	UID int `json:"uid" binding:"required"`
}

// MultiUIDForm 像AddInterviewers，Pass，Out， Enroll，DeleteInterviewers 都需要 uid 数组
type MultiUIDForm struct {
	UID  []int `json:"uids" binding:"required"`
	Type int   `json:"type"` // 1表示第一轮，2表示第二轮，3表示录取，4表示淘汰
}

type Time struct {
	Start int64 `binding:"required" json:"start"`
	End   int64 `binding:"required" json:"end"`
}
