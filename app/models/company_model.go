package models

import "time"

type Company struct {
	Company_Id     uint      `db:"company_id"      json:"company_id"`
	CreatedAt      time.Time `db:"created_at"      json:"created_at"`
	Company_status string    `db:"company_status"  json:"company_status"`
	Company_name   string    `db:"company_name"    json:"company_name"`
	Inn            string    `db:"inn"             json:"inn"`
	Kpp            string    `db:"kpp"             json:"kpp"`
	Ogrn           string    `db:"ogrn"            json:"ogrn"`
	Ur_adress      string    `db:"ur_adress"       json:"ur_adress"`
	Mail_adress    string    `db:"mail_adress"     json:"mail_adress"`
	Phone          string    `db:"phone"           json:"phone"`
	Email          string    `db:"email"           json:"email"`
	Director       string    `db:"director"        json:"director"`
}
