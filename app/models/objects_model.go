package models

import "time"

type Objects struct {
	Object_Id uint      `db:"object_id"       json:"object_id"`
	CreatedAt time.Time `db:"created_at"      json:"created_at"`
	Typereal  string    `db:"typereal"        json:"typereal"`
	City      string    `db:"city"            json:"city"`
	House     string    `db:"house"           json:"house"`
	Flat      string    `db:"flat"            json:"flat"`
	Square    string    `db:"square"          json:"square"`
}
