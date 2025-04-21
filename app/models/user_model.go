package models

import "time"

// UserType определяет структуру пользователя
type UserType struct {
	UserID            uint      `db:"user_id"          json:"userId,omitempty"`
	CreatedUser       time.Time `db:"created_user"     json:"createdUser,omitempty"`
	BdUsed            string    `db:"bd_used"          json:"bdUsed,omitempty"`
	FirstName         string    `db:"first_name"       json:"firstName,omitempty"`
	PatronymicName    string    `db:"patronymic_name"  json:"patronymicName,omitempty"`
	LastName          string    `db:"last_name"        json:"lastName,omitempty"`
	UserEmail         string    `db:"user_email"       json:"userEmail"`
	EmailVerification bool      `db:"email_verification"       json:"emailVerification"`
	UserPhone         string    `db:"user_phone"       json:"userPhone,omitempty"`
	Password          string    `db:"password"         json:"password,omitempty"`
}

// UserRights определяет права пользователя
type UserRights struct {
	RightsID      uint      `db:"rights_id"      json:"rightsId,omitempty"`
	UserID        uint      `db:"user_id"        json:"userId,omitempty"`
	CreatedRights time.Time `db:"created_rights" json:"createdRights,omitempty"`
	UserBD        string    `db:"user_bd"        json:"userBd,omitempty"`
	Holding       string    `db:"holding"        json:"holding,omitempty"`
	UserRole      string    `db:"user_role"      json:"userRole,omitempty"`
}

// UserRole определяет роли пользователя
type UserRole string

const (
	RoleMember   UserRole = "member"
	RoleAdmin    UserRole = "admin"
	RoleDirector UserRole = "director"
	RoleManager  UserRole = "manager"
)
