package models

import "time"

type BancStatementType struct {
	StatementId      uint      `db:"statement_id"        json:"statementId,omitempty"`
	CreatedStatement time.Time `db:"created_statement"   json:"createdStatement,omitempty"`
	RSchetId         string    `db:"r_schet_id"          json:"rSchetId,omitempty"`
	RSchet           string    `db:"r_schet"             json:"rSchet,omitempty"`
	DateTransaction  time.Time `db:"date_transaction"    json:"dateTransaction,omitempty"`
	CompanyId        uint      `db:"company_id"          json:"companyId,omitempty"`
	ContragentId     uint      `db:"contragent_id"       json:"contragentId,omitempty"`
	BalanceBeginDay  string    `db:"balance_begin_day"   json:"balanceBeginDay,omitempty"`
	Kredit           string    `db:"kredit"              json:"kredit,omitempty"`
	Debit            string    `db:"debit"               json:"debit,omitempty"`
	BalanceEndDay    string    `db:"balance_end_day"     json:"balanceEndDay,omitempty"`
	BasisPayment     string    `db:"basis_payment"       json:"basisPayment,omitempty"`
	AuthorId         string    `db:"author_id"           json:"authorId,omitempty"`
}
