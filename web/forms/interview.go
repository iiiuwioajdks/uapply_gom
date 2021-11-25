package forms

type Interviewer struct {
	UID int `json:"uid" binding:"required,min=1"`
}
