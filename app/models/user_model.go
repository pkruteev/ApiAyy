package models

import "time"

// UserType определяет структуру пользователя
type UserType struct {
	UserID         uint      `db:"user_id"               json:"user_id,omitempty"`
	CreatedUser    time.Time `db:"created_user"          json:"created_user,omitempty"`
	BdUsed         uint      `db:"bd_used"               json:"bd_used,omitempty"`
	Holding        string    `db:"holding"               json:"holding,omitempty"`
	FirstName      string    `db:"first_name"            json:"first_name,omitempty"`
	PatronymicName string    `db:"patronymic_name"       json:"patronymic_name,omitempty"`
	LastName       string    `db:"last_name"             json:"last_name,omitempty"`
	RightsID       uint      `db:"rights_id"             json:"rights_id,omitempty"`
	UserBD         uint      `db:"user_bd"               json:"user_bd,omitempty"`
	UserRights     string    `db:"user_rights"           json:"user_rights,omitempty"`
	UserEmail      string    `db:"user_email"            json:"user_email"`
	UserPhone      string    `db:"user_phone"            json:"user_phone,omitempty"`
	Password       string    `db:"password"              json:"password,omitempty"`
}

// UserRole определяет роли пользователя
type UserRigh struct {
	Member   string `db:"member" json:"member,omitempty"`
	Admin    string `db:"admin" json:"admin,omitempty"`
	Director string `db:"director" json:"director,omitempty"`
	Manager  string `db:"manager" json:"manager,omitempty"`
}

type UserResponse struct {
	UserID         uint   `db:"user_id"          json:"user_id,omitempty"`
	BdUsed         uint   `db:"bd_used"               json:"bd_used,omitempty"`
	FirstName      string `db:"first_name"       json:"first_name,omitempty"`
	LastName       string `db:"last_name"        json:"last_name,omitempty"`
	PatronymicName string `db:"patronymic_name"  json:"patronymic_name,omitempty"`
	UserRights     string `db:"user_rights"      json:"user_rights,omitempty"`
	UserEmail      string `db:"user_email"       json:"user_email"`
	UserPhone      string `db:"user_phone"       json:"user_phone,omitempty"`
}

type UserResponses struct {
	UserID         uint   `db:"user_id" json:"user_id,omitempty"`
	FirstName      string `db:"first_name" json:"first_name,omitempty"`
	LastName       string `db:"last_name" json:"last_name,omitempty"`
	PatronymicName string `db:"patronymic_name" json:"patronymic_name,omitempty"`
	UserRights     string `db:"user_rights" json:"user_rights,omitempty"`
	UserEmail      string `db:"user_email" json:"user_email"`
	UserPhone      string `db:"user_phone" json:"user_phone,omitempty"`
}
