package models

import "time"

type Company struct {
	CompanyId      uint      `db:"company_id"      json:"companyId,omitempty"`
	CreateCompany  time.Time `db:"create_company"  json:"createCompany,omitempty"`
	IsCounterparty bool      `db:"is_counterparty" json:"isCounterparty"`
	Status         string    `db:"status"          json:"status"`
	Name           string    `db:"name"            json:"name"`
	Inn            string    `db:"inn"             json:"inn"`
	Kpp            string    `db:"kpp"             json:"kpp"`
	Ogrn           string    `db:"ogrn"            json:"ogrn"`
	DataOgrn       time.Time `db:"data_ogrn"       json:"dataOgrn,omitempty"`
	Ogrnip         string    `db:"ogrnip"          json:"ogrnip"`
	DataOgrnip     time.Time `db:"data_ogrnip"     json:"dataOgrnip,omitempty"`
	UrAddress      string    `db:"ur_address"       json:"urAddress"`
	MailAddress    string    `db:"mail_address"     json:"mailAddress"`
	Phone          string    `db:"phone"           json:"phone,omitempty"`
	Email          string    `db:"email"           json:"email"`
	Director       string    `db:"director"        json:"director"`
}

type CompanyStatus struct {
	OOO   string `db:"ooo"      json:"ooo,omitempty"`
	Ip    string `db:"ip"       json:"ip,omitempty"`
	Fizik string `db:"fiz_lico" json:"fizLico,omitempty"`
	Samo  string `db:"samo"     json:"samo,omitempty"`
}
