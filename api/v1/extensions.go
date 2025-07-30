package v1

import (
	"fmt"
	"regexp"
	"time"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
	"github.com/go-playground/validator/v10"
	"github.com/oapi-codegen/runtime/types"
)

// Entity converts a v1.TransactionForm to an entity.Transaction for business logic processing.
// This method transforms the HTTP request data into the internal domain model representation,
// parsing the date and creating a properly formatted transaction entity.
func (t TransactionForm) Entity() (entity.Transaction, error) {
	date, _ := time.Parse(t.Date.String(), types.DateFormat)
	te := entity.NewTransaction(entity.TransactionKind(t.Kind), t.Account, date, t.Amount, t.Content)

	return *te, nil
}

// CreateTransactionFormValidation provides custom validation logic for v1.TransactionForm structures.
// It implements the validator.StructLevel interface for complex validation rules.
func CreateTransactionFormValidation(sl validator.StructLevel) {
	form := sl.Current().Interface().(TransactionForm)

	if form.Kind != Credit && form.Kind != Debit {
		sl.ReportError(form.Kind, "kind", "kind", "should be credit or debit", "")
	}

	// Date validation is handled by the openapi_types.Date type itself
}

// UpdateTransactionInfoFormValidation provides custom validation logic for v1.TransactionInfoForm structures.
// It implements the validator.StructLevel interface for validation rules.
func UpdateTransactionInfoFormValidation(sl validator.StructLevel) {
	form := sl.Current().Interface().(TransactionInfoForm)

	// Info field can be empty or up to a reasonable length
	if form.Info != nil && len(*form.Info) > 1000 {
		sl.ReportError(form.Info, "info", "info", "max length is 1000 characters", "")
	}
}

// Entity converts a v1.RuleForm to an entity.Rule for business logic processing.
// This method transforms the HTTP request data into the internal domain model representation,
// converting the labels map into a slice of entity.Label structures.
func (r RuleForm) Entity() entity.Rule {
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

	if form.Name == "" || len(form.Name) > 255 {
		sl.ReportError(form.Name, "name", "name", "lt 255", "")
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

// Entity converts a v1.LabelForm to an entity.Label for business logic processing.
// This method transforms the HTTP request data into the internal domain model representation.
// Note: ID and CreatedAt are not set as they will be populated by the service layer.
func (l LabelForm) Entity() entity.Label {
	return entity.Label{
		Key:   l.Key,
		Value: l.Value,
	}
}

// NewTransaction converts an entity.Transaction to a v1.Transaction for API responses.
// This method transforms the internal domain model into the HTTP response format,
// including proper HREF generation and date formatting.
func NewTransaction(t entity.Transaction) Transaction {
	// Convert labels
	labels := make([]TransactionLabelAssociation, 0, len(t.Labels))
	for _, a := range t.Labels {
		href := fmt.Sprintf("/api/v1/labels/%d", a.Label.ID)
		key := a.Label.Key
		value := a.Label.Value

		tLabelAssociation := TransactionLabelAssociation{
			Href:  &href,
			Key:   &key,
			Value: &value,
		}
		if a.RuleID != nil {
			ruleHref := fmt.Sprintf("/api/v1/rules/%s", *a.RuleID)
			tLabelAssociation.RuleHref = &ruleHref
		}
		labels = append(labels, tLabelAssociation)
	}

	// Convert main transaction fields
	href := fmt.Sprintf("/api/v1/transactions/%d", t.ID)
	account := t.Account
	amount := t.Amount
	date := t.Date
	description := t.RawContent
	kind := string(t.Kind)

	return Transaction{
		Account:     &account,
		Amount:      &amount,
		Date:        &date,
		Description: &description,
		Href:        &href,
		Info:        t.Info,
		Kind:        &kind,
		Labels:      &labels,
	}
}

// FromEntityRule converts an entity.Rule to a v1.Rule for API responses.
// This method transforms the internal domain model into the HTTP response format.
func FromEntityRule(rule entity.Rule) Rule {
	href := fmt.Sprintf("/api/v1/rules/%s", rule.Name)
	pattern := rule.Pattern

	labels := make([]Label, 0, len(rule.Labels))
	for _, label := range rule.Labels {
		labelHref := fmt.Sprintf("/api/v1/labels/%d", label.ID)
		labels = append(labels, Label{
			Href:  labelHref,
			Key:   label.Key,
			Value: label.Value,
		})
	}

	return Rule{
		Href:    href,
		Labels:  &labels,
		Name:    rule.Name,
		Pattern: &pattern,
	}
}

// NewRules converts a slice of entity.Rule to a v1.Rules collection for API responses.
func NewRules(rules []entity.Rule) Rules {
	r := Rules{
		Rules: make([]Rule, 0, len(rules)),
		Total: len(rules),
	}
	for _, rr := range rules {
		r.Rules = append(r.Rules, FromEntityRule(rr))
	}
	return r
}

// FromEntityLabel converts an entity.Label to a v1.Label for API responses.
// This method transforms the internal domain model into the HTTP response format.
func FromEntityLabel(l entity.Label) Label {
	href := fmt.Sprintf("/api/v1/labels/%d", l.ID)

	rules := make([]Rule, 0, len(l.Rules))
	for _, rule := range l.Rules {
		ruleHref := fmt.Sprintf("/api/v1/rules/%s", rule)
		rules = append(rules, Rule{
			Href: ruleHref,
			Name: rule,
		})
	}

	return Label{
		Href:  href,
		Key:   l.Key,
		Value: l.Value,
		Rules: &rules,
	}
}

// NewLabels converts a slice of entity.Label to a v1.Labels collection for API responses.
func NewLabels(labels []entity.Label) Labels {
	mlabels := Labels{
		Labels: make([]Label, 0, len(labels)),
		Total:  len(labels),
	}

	for _, label := range labels {
		mlabels.Labels = append(mlabels.Labels, FromEntityLabel(label))
	}

	return mlabels
}

// NewTransactions creates a new v1.Transactions response structure with the given
// total count, time range, and transaction data.
func NewTransactions(total int, start, end time.Time) Transactions {
	startDate := types.Date{Time: start}
	endDate := types.Date{Time: end}
	totalPtr := total
	items := make([]Transaction, 0, total)

	return Transactions{
		Start: &startDate,
		End:   &endDate,
		Total: &totalPtr,
		Items: &items,
	}
}
