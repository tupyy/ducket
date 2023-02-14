package entity

import "regexp"

type Rule struct {
	Pattern regexp.Regexp
	Labels  []Label
	Fields  map[string]string
}
