package models

import "time"

type Blog struct {
	Id        string    `json:"id" db:"Id"`
	UserId    string    `json:"user_id" db:"UserId"`
	Title     string    `json:"title" db:"Title"`
	Content   string    `json:"content" db:"Content"`
	Category  string    `json:"category" db:"Category"`
	Tags      []string  `json:"tags" db:"Tags"`
	CreatedAt time.Time `json:"created_at" db:"CreatedAt"`
	UpdatedAt time.Time `json:"updated_at" db:"UpdatedAt"`
}
