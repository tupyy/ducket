package entity

import "time"

type Rule struct {
	ID        int
	Name      string
	Filter    string
	Tags      []string
	CreatedAt time.Time
}
