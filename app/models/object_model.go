package models

import "time"

type Objects struct {
	ObjectId  uint      `db:"object_id"       json:"objectId"`
	CreatedOb time.Time `db:"created_ob"      json:"createdOb"`
	TypeReal  string    `db:"typereal"        json:"typeReal"`
	City      string    `db:"city"            json:"city"`
	Street    string    `db:"street"          json:"street"`
	House     string    `db:"house"           json:"house"`
	Flat      string    `db:"flat"            json:"flat"`
	Square    string    `db:"square"          json:"square"`
	Floor     string    `db:"floor"           json:"floor"`
}
