package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type StatementType struct {
	StatementId      uint            `db:"statement_id"        json:"statementId,omitempty"`
	CreatedStatement time.Time       `db:"created_statement"   json:"createdStatement,omitempty"`
	RSchetId         string          `db:"r_schet_id"          json:"rSchetId,omitempty"`
	RSchet           string          `db:"r_schet"             json:"rSchet,omitempty"`
	DateTransaction  time.Time       `db:"date_transaction"    json:"dateTransaction,omitempty"`
	CompanyId        uint            `db:"company_id"          json:"companyId,omitempty"`
	ContragentId     uint            `db:"contragent_id"       json:"contragentId,omitempty"`
	BalanceBeginDay  decimal.Decimal `db:"balance_begin_day"   json:"balanceBeginDay,omitempty"`
	Kredit           decimal.Decimal `db:"kredit"              json:"kredit,omitempty"`
	Debit            decimal.Decimal `db:"debit"               json:"debit,omitempty"`
	BalanceEndDay    decimal.Decimal `db:"balance_end_day"     json:"balanceEndDay,omitempty"`
	BasisPayment     string          `db:"basis_payment"       json:"basisPayment,omitempty"`
	AuthorId         string          `db:"author_id"           json:"authorId,omitempty"`
}
