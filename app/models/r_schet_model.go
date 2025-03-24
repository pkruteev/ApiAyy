package models

import "time"

type RSchet struct {
	RSchetId     uint      `db:"r_schet_id"      json:"rSchetId"`
	CreatedSchet time.Time `db:"created_schet"   json:"createdSchet"`
	CompanyId    string    `db:"company_id"      json:"companyId"`
	RSchet       string    `db:"r_schet"         json:"rSchet"`
	BankName     string    `db:"bank_name"       json:"bankName"`
	BankBic      string    `db:"bank_bic"        json:"bankBic"`
	KorSchet     string    `db:"kor_schet"       json:"korSchet"`
}
