package inbound

type TagForm struct {
	Value string `form:"value" json:"value" binding:"required" validate:"required,min=1,max=20"`
}
