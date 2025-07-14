package inbound

import (
	"regexp"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
	"github.com/go-playground/validator/v10"
)

type RuleForm struct {
	Name    string            `form:"name" json:"name" binding:"required"`
	Pattern string            `form:"pattern" json:"pattern" binding:"required"`
	Labels  map[string]string `form:"labels" json:"labels" binding:"required"`
}

// FormToEntity converts a RuleForm to an entity.Rule for business logic processing.
func FormToEntity(r RuleForm) entity.Rule {
	labels := make([]entity.Label, 0, len(r.Labels))
	for key, value := range r.Labels {
		labels = append(labels, entity.Label{
			Key:   key,
			Value: value,
		})
	}
	return entity.NewRule(r.Name, r.Pattern, labels...)
}

// RuleFormValidation provides custom validation logic for RuleForm structures.
// It implements the validator.StructLevel interface for complex validation rules.
func RuleFormValidation(sl validator.StructLevel) {
	form := sl.Current().Interface().(RuleForm)

	if form.Name == "" || len(form.Name) > 20 {
		sl.ReportError(form.Name, "name", "name", "lt 20", "")
	}

	if _, err := regexp.Compile(form.Pattern); err != nil {
		sl.ReportError(form.Pattern, "pattern", "pattern", "must compile", "")
	}

	if len(form.Labels) == 0 {
		sl.ReportError(form.Labels, "labels", "labels", "ge 0", "")
	}

	for _, l := range form.Labels {
		if len(l) > 20 {
			sl.ReportError(form.Labels, "label", "label", "lt 20", "")
		}
	}
}

type UpdateRuleForm struct {
	Pattern string            `form:"pattern" json:"pattern" binding:"required"`
	Labels  map[string]string `form:"labels" json:"labels" binding:"required"`
}

// UpdateRuleFormValidation provides custom validation logic for UpdateRuleForm structures.
// It implements the validator.StructLevel interface for update-specific validation rules.
func UpdateRuleFormValidation(sl validator.StructLevel) {
	form := sl.Current().Interface().(UpdateRuleForm)

	if _, err := regexp.Compile(form.Pattern); err != nil {
		sl.ReportError(form.Pattern, "pattern", "pattern", "must compile", "")
	}

	if len(form.Labels) == 0 {
		sl.ReportError(form.Labels, "labels", "labels", "ge 0", "")
	}

	for _, l := range form.Labels {
		if len(l) > 20 {
			sl.ReportError(form.Labels, "label", "label", "lt 20", "")
		}
	}
}
