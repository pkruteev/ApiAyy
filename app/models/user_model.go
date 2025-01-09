package models

import "time"

// UserType определяет структуру пользователя
type UserType struct {
	User_ID        uint      `db:"user_id"          json:"user_id,omitempty"`
	CreatedAt      time.Time `db:"created_at"       json:"created_at,omitempty"`
	FirstName      string    `db:"first_name"       json:"first_name,omitempty"`
	LastName       string    `db:"last_name"        json:"last_name,omitempty"`
	PatronymicName string    `db:"patronymic_name"  json:"patronymic_name,omitempty"`
	Rights_Id      uint      `db:"rights_id"        json:"rights_id,omitempty"`
	User_Company   string    `db:"user_company"     json:"user_company,omitempty"`
	User_Right     string    `db:"user_right"       json:"user_right,omitempty"`
	UserEmail      string    `db:"user_email"       json:"user_email"`
	UserPhone      string    `db:"user_phone"       json:"user_phone,omitempty"`
	Password       string    `db:"password"         json:"password"`
}

// UserRole определяет роли пользователя
type UserRights struct {
	Member   string `db:"member" json:"member,omitempty"`
	Admin    string `db:"admin" json:"admin,omitempty"`
	Director string `db:"director" json:"director,omitempty"`
	Manager  string `db:"manager" json:"manager,omitempty"`
}
