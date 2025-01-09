package models

import "time"

type Contragency struct {
	Contragency_Id     uint      `db:"contragency_id"      json:"contragency_id"`
	CreatedAt          time.Time `db:"created_at"          json:"created_at"`
	Contragency_status string    `db:"contragency_status"  json:"contragency_status"`
	Contragency_name   string    `db:"contragency_name"    json:"contragency_name"`
	Inn                string    `db:"inn"                 json:"inn"`
	Kpp                string    `db:"kpp"                 json:"kpp"`
	Ogrn               string    `db:"ogrn"                json:"ogrn"`
	Ur_adress          string    `db:"ur_adress"           json:"ur_adress"`
	Mail_adress        string    `db:"mail_adress"         json:"mail_adress"`
	Phone              string    `db:"phone"               json:"phone"`
	Email              string    `db:"email"               json:"email"`
	Director           string    `db:"director"            json:"director"`
}
