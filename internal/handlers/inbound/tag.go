package inbound

import (
	"github.com/go-playground/validator/v10"
)

type TagForm struct {
	Value string `form:"value" json:"value" binding:"required" validate:"required"`
}

func TagFormValidation(sl validator.StructLevel) {
	form := sl.Current().Interface().(TagForm)

	if len(form.Value) > 20 {
		sl.ReportError(form.Value, "value", "Value", "gt 0 and lt 20", "")
	}
}
