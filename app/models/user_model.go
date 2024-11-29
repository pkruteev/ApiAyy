package models

// UserType определяет структуру пользователя
type UserType struct {
	UserID         int      `db:"user_id" json:"user_id"`
	FirstName      string   `db:"first_name" json:"first_name"`
	PatronymicName string   `db:"patronymic_name" json:"patronymic_name,omitempty"`
	LastName       string   `db:"last_name" json:"last_name,omitempty"`
	UserEmail      string   `db:"user_email" json:"user_email"`
	UserPhone      string   `db:"user_phone" json:"user_phone,omitempty"`
	UserCompany    []string `db:"user_company" json:"user_company,omitempty"`
	Password       string   `db:"password" json:"password,omitempty"`
	UserRole       UserRole `db:"user_role" json:"user_role" validate:"required,dive"`
}

// UserRole определяет роли пользователя
type UserRole struct {
	Member   string `json:"member"`
	Admin    string `json:"admin"`
	Director string `json:"director"`
	Manager  string `json:"manager"`
}
