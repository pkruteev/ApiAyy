package models

import "time"

type Contract struct {
	ContractId      uint      `db:"contract_id"            json:"contractId"`
	CreatedContract time.Time `db:"created_contract"       json:"createdContract"`
	ContractNumber  string    `db:"contract_number"        json:"contractNumber"`
	DateStart       time.Time `db:"date_start"             json:"dateStart"`
	DateEnd         time.Time `db:"date_end"               json:"dateEnd"`
	DateStartRent   time.Time `db:"date_start_rent"        json:"dateStartRent"`
	DateEndRent     time.Time `db:"date_end_rent"          json:"dateEndRent"`
	ObjectId        uint      `db:"object_id"              json:"objectId"`
	CompanyId       uint      `db:"company_id"             json:"companyId"`
	CounterpartyId  uint      `db:"counterparty_id"        json:"counterpartyId"`
	RentPay         string    `db:"rent_pay"               json:"rentPay"`
	DayPaymentRent  uint      `db:"day_payment_rent"       json:"dayPaymentRent"`
	RentPrePay      string    `db:"rent_pre_pay"           json:"rentPrePay"`
	DateRentPrePay  time.Time `db:"date_rent_prepay"       json:"dateRentPrePay"`
	UtilityBills    bool      `db:"utility_bills"          json:"utilityBills"`
}
