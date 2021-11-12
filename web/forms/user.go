package forms

// UserInfoReq 用户简历请求，为了能在保存草稿时复用，此处不做 binding 要求
type UserInfoReq struct {
	UID     int32  `json:"uid"`
	Name    string `json:"name"`
	StuNum  string `json:"stu_num"`
	Address string `json:"address"`
	Major   string `json:"major"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Sex     int8   `json:"sex"`
	Intro   string `json:"intro"`
}
