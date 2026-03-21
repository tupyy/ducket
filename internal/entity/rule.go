package entity

import "time"

type Rule struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Filter    string    `json:"filter"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
}
