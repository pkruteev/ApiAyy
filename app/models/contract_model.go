package models

import "time"

type ContractType struct {
	ContractId            uint      `db:"contract_id"             json:"contractId"`
	CreatedContract       time.Time `db:"created_contract"        json:"createdContract,omitempty"`
	DateSigningContract   time.Time `db:"date_signing_contract"   json:"dateSigningContract"`
	ContractNumber        string    `db:"contract_number"         json:"contractNumber"`
	DateStart             time.Time `db:"date_start"              json:"dateStart"`
	DateEnd               time.Time `db:"date_end"                json:"dateEnd"`
	DateStartRent         time.Time `db:"date_start_rent"         json:"dateStartRent"`
	DateEndRent           time.Time `db:"date_end_rent"           json:"dateEndRent"`
	ObjectId              uint      `db:"object_id"               json:"objectId"`
	CompanyId             uint      `db:"company_id"              json:"companyId"`
	RSchetId              uint      `db:"r_schet_id"              json:"RSchetId"`
	PaymentMethodBanc     bool      `db:"payment_method_banc"     json:"paymentMethodBanc,omitempty"`
	PaymentMethodCash     bool      `db:"payment_method_cash"     json:"paymentMethodCash,omitempty"`
	CounterpartyId        uint      `db:"counterparty_id"         json:"counterpartyId"`
	RentPay               string    `db:"rent_pay"                json:"rentPay"`
	DayPaymentRent        string    `db:"day_payment_rent"        json:"dayPaymentRent"`
	RentPrePay            string    `db:"rent_pre_pay"            json:"rentPrePay"`
	DateRentPrePay        time.Time `db:"date_rent_pre_pay"       json:"dateRentPrePay"`
	IsUtilitiesIncluded   bool      `db:"is_utilities_included"   json:"isUtilitiesIncluded"`
	IsWaterIncluded       bool      `db:"is_utilities_included"   json:"isWaterIncluded"`
	IsConciergeIncluded   bool      `db:"is_concierge_included"   json:"isConciergeIncluded"`
	IsElectricityIncluded bool      `db:"is_electricity_included" json:"isElectricityIncluded"`
	TypeReal              string    `db:"type_real"               json:"typeReal"` // "commercial" or "residential"
	TerminationDate       time.Time `db:"termination_date"        json:"terminationDate"`
}
