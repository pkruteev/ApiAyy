package models

// UserType определяет структуру пользователя
type MyUserType struct {
	RightsID       uint   `db:"rights_id"        json:"rightsId,omitempty"`
	UserID         uint   `db:"user_id"          json:"userId,omitempty"`
	FirstName      string `db:"first_name"       json:"firstName,omitempty"`
	PatronymicName string `db:"patronymic_name"  json:"patronymicName,omitempty"`
	LastName       string `db:"last_name"        json:"lastName,omitempty"`
	UserEmail      string `db:"user_email"       json:"userEmail"`
	UserPhone      string `db:"user_phone"       json:"userPhone,omitempty"`
	UserRole       string `db:"user_role"        json:"userRole,omitempty"`
}
