package inbound

import (
	"regexp"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
	"github.com/go-playground/validator/v10"
)

type RuleForm struct {
	Name    string   `form:"name" json:"name" binding:"required"`
	Pattern string   `form:"pattern" json:"pattern" binding:"required"`
	Tags    []string `form:"tags" json:"tags" binding:"required"`
}

func FormToEntity(r RuleForm) entity.Rule {
	return entity.NewRule("id", r.Name, r.Pattern, r.Tags...)
}

func RuleFormValidation(sl validator.StructLevel) {
	form := sl.Current().Interface().(RuleForm)

	if len(form.Name) > 20 {
		sl.ReportError(form.Name, "name", "name", "lt 20", "")
	}

	if _, err := regexp.Compile(form.Pattern); err != nil {
		sl.ReportError(form.Pattern, "pattern", "pattern", "must compile", "")
	}

	if len(form.Tags) == 0 {
		sl.ReportError(form.Tags, "tags", "tags", "ge 0", "")
	}

	for _, t := range form.Tags {
		if len(t) > 20 {
			sl.ReportError(form.Tags, "tag", "tag", "lt 20", "")
		}
	}
}
