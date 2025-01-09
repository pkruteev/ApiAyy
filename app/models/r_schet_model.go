package models

import "time"

type R_schets struct {
	R_schet_Id uint      `db:"r_schet_id"      json:"r_schet_id"`
	CreatedAt  time.Time `db:"created_at"      json:"created_at"`
	R_schet    string    `db:"r_schet"         json:"r_schet"`
	Bank_name  string    `db:"bank_name"       json:"bank_name"`
	Bank_bic   string    `db:"bank_bic"        json:"bank_bic"`
	Kor_schet  string    `db:"kor_schet"       json:"kor_schet"`
}
