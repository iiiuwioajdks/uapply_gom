package forms

// UserInfoReq 用户简历请求
type UserInfoReq struct {
	UID     int32  `json:"uid"`
	Name    string `json:"name" binding:"required"`
	StuNum  string `json:"stu_num" binding:"required"`
	Address string `json:"address,omitempty"`
	Major   string `json:"major,omitempty"`
	Phone   string `json:"phone" binding:"required,mobile"`
	Email   string `json:"email" binding:"required,email"`
	Sex     int8   `json:"sex" binding:"required,oneof=1 2"` // 1为男，2为女
	Intro   string `json:"intro" binding:"required,max=100"`
}

// UserRegisterInfo 用户报名表单
type UserRegisterInfo struct {
	UID            int32 `json:"uid"`
	OrganizationID int   `json:"organization_id" binding:"required"`
	DepartmentID   int   `json:"department_id" binding:"required"`
}

// UserResumeInfo 更新简历和保存到草稿箱的表单
type UserResumeInfo struct {
	UID     int32  `json:"uid"`
	Name    string `json:"name"`
	StuNum  string `json:"stu_num"`
	Address string `json:"address"`
	Major   string `json:"major"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Sex     int8   `json:"sex"` // 1为男，2为女
	Intro   string `json:"intro" binding:"min=0,max=100"`
}
