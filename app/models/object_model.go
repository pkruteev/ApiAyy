package models

import "time"

type Objects struct {
	ObjectId      uint      `db:"object_id"       json:"objectId"`
	CompanyId     uint      `db:"company_id"      json:"companyId,omitempty"`
	CadastrNumber string    `db:"cadastr_number"  json:"cadastrNumber,omitempty"`
	CreatedOb     time.Time `db:"created_ob"      json:"createdOb"`
	TypeReal      string    `db:"typereal"        json:"typeReal"`
	City          string    `db:"city"            json:"city"`
	Street        string    `db:"street"          json:"street"`
	House         string    `db:"house"           json:"house"`
	Flat          string    `db:"flat"            json:"flat"`
	Square        string    `db:"square"          json:"square"`
	Floor         string    `db:"floor"           json:"floor"`
}
