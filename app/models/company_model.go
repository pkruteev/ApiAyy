package models

import "time"

type Company struct {
	Id_company     uint      `db:"id_company"      json:"id_company"`
	CreatedAt      time.Time `db:"created_at"      json:"created_at"`
	Status_company string    `db:"status_company"  json:"status"`
	Company_name   string    `db:"company_name"    json:"company_name"`
	Inn            string    `db:"inn"             json:"inn"`
	Kpp            string    `db:"kpp"             json:"kpp"`
	Ogrn           string    `db:"ogrn"            json:"ogrn"`
	Ur_adress      string    `db:"ur_adress"       json:"ur_adress"`
	Mail_adress    string    `db:"mail_adress"     json:"mail_adress"`
	Bank_name      string    `db:"bank_name"       json:"bank_name"`
	Bank_bic       string    `db:"bank_bic"        json:"bank_bic"`
	Kor_schet      string    `db:"kor_schet"       json:"kor_schet"`
	R_schet        string    `db:"r_schet"         json:"r_schet"`
	Phone          string    `db:"phone"           json:"phone"`
	Email          string    `db:"email"           json:"email"`
	Director       string    `db:"director"        json:"director"`
}
