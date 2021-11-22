package response

type IntervieweeRsp struct {
	UID            int    `json:"uid"`
	Name           string `json:"name"`
	StuNum         string `json:"stu_num"`
	Address        string `json:"address"`
	Major          string `json:"major"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	Sex            int8   `json:"sex"`
	Intro          string `json:"intro"`
	OrganizationID int    `json:"organization_id"`
	DepartmentID   int    `json:"department_id"`
	FirstStatus    int8   `json:"first_status"`
	SecondStatus   int8   `json:"second_status"`
	FinalStatus    int8   `json:"final_status"`
}