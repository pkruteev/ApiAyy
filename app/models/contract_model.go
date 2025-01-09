package models

import "time"

type Contracts struct {
	Contract_Id         uint      `db:"contract_id"            json:"contract_id"`
	CreatedAt           time.Time `db:"created_at"             json:"created_at"`
	Contract_number     string    `db:"contract_number"        json:"contract_number"`
	Data_start_contract time.Time `db:"data_start_contract"    json:"data_start_contract"`
	Data_end_contract   time.Time `db:"data_end_contract"      json:"data_end_contract"`
	Data_start_rent     time.Time `db:"data_start_rent"        json:"data_start_rent"`
	Data_end_rent       time.Time `db:"data_end_rent"          json:"data_end_rent"`
	Object_id           string    `db:"object_id"              json:"object_id"`
	Company_id          string    `db:"company_id"             json:"company_id"`
	Contragency_id      string    `db:"contragency_id"         json:"contragency_id"`
	Service_contract    bool      `db:"service_contract"       json:"service_contract"`
	Rent_pay            string    `db:"rent_pay"               json:"rent_pay"`
	Rent_pre_pay        string    `db:"rent_pre_pay"           json:"rent_pre_pay"`
}
