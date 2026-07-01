package models

import "time"

type User struct {
	Id        string    `json:"id" db:"Id"`
	Name      string    `json:"name" db:"Name"`
	Email     string    `json:"email" db:"Email"`
	Password  string    `json:"password" db:"Password"`
	Role      string    `json:"role" db:"Role"`
	CreatedAt time.Time `json:"created_at" db:"CreatedAt"`
	UpdatedAt time.Time `json:"updated_at" db:"UpdatedAt"`
}
