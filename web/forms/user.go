package forms

// UserInfoReq 用户简历请求，为了能在保存草稿时复用，此处不做 binding 要求
type UserInfoReq struct {
	UID     int32  `json:"uid"`
	Name    string `json:"name" binding:"required"`
	StuNum  string `json:"stu_num" binding:"required"`
	Address string `json:"address"`
	Major   string `json:"major"`
	Phone   string `json:"phone" binding:"required,mobile"`
	Email   string `json:"email" binding:"required,email"`
	Sex     int8   `json:"sex" binding:"required"` // 1为男，2为女
	Intro   string `json:"intro" binding:"required,max=100"`
}
