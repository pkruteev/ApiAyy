package queries

import (
	"ApiAyy/app/models"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ContractQueries struct {
	*sqlx.DB
}

// GetContracts returns all contracts
func (q *ContractQueries) GetContracts() ([]models.ContractType, error) {
	contracts := []models.ContractType{}
	query := `SELECT * FROM contracts ORDER BY contract_id`
	err := q.Select(&contracts, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get contracts: %w", err)
	}
	return contracts, nil
}

// GetContract returns a single contract by ID
func (q *ContractQueries) GetContract(id uint) (models.ContractType, error) {
	var contract models.ContractType
	query := `SELECT * FROM contracts WHERE contract_id = $1 LIMIT 1`
	err := q.Get(&contract, query, id)
	if err != nil {
		return models.ContractType{}, fmt.Errorf("failed to get contract: %w", err)
	}
	return contract, nil
}

// CreateContract creates a new contract
func (q *ContractQueries) CreateContract(c *models.ContractType) error {
	query := `INSERT INTO contracts (
		created_contract, date_signing_contract, contract_number, date_start, date_end, 
		date_start_rent, date_end_rent, object_id, company_id, r_schet_id,
		payment_method_banc, payment_method_cash, counterparty_id, rent_pay,
		day_payment_rent, rent_pre_pay, date_rent_pre_pay, 
		is_utilities_included, is_water_included,
		is_concierge_included, is_electricity_included, type_real,
		termination_date
	) VALUES (
		:created_contract, :date_signing_contract, :contract_number, :date_start, :date_end, 
		:date_start_rent, :date_end_rent, :object_id, :company_id, :r_schet_id,
		:payment_method_banc, :payment_method_cash, :counterparty_id, :rent_pay,
		:day_payment_rent, :rent_pre_pay, :date_rent_pre_pay, 
		:is_utilities_included, :is_water_included,
		:is_concierge_included, :is_electricity_included, :type_real,
		:termination_date
	) RETURNING contract_id`

	_, err := q.NamedExec(query, c)
	if err != nil {
		return fmt.Errorf("failed to create contract: %w", err)
	}
	return nil
}

// UpdateContract updates contract data
func (q *ContractQueries) UpdateContract(id uint, c *models.ContractType) error {
	c.ContractId = id // Set ID from parameter

	query := `UPDATE contracts SET 
		created_contract = :created_contract,
		date_signing_contract = :date_signing_contract, 
		contract_number = :contract_number, 
		date_start = :date_start, 
		date_end = :date_end, 
		date_start_rent = :date_start_rent, 
		date_end_rent = :date_end_rent, 
		object_id = :object_id, 
		company_id = :company_id, 
		r_schet_id = :r_schet_id,
		payment_method_banc = :payment_method_banc, 
		payment_method_cash = :payment_method_cash, 
		counterparty_id = :counterparty_id, 
		rent_pay = :rent_pay,
		day_payment_rent = :day_payment_rent, 
		rent_pre_pay = :rent_pre_pay, 
		date_rent_pre_pay = :date_rent_pre_pay, 
		is_utilities_included = :is_utilities_included, 
		is_water_included = :is_water_included,
		is_concierge_included = :is_concierge_included, 
		is_electricity_included = :is_electricity_included, 
		type_real = :type_real,
		termination_date = :termination_date
	WHERE contract_id = :contract_id`

	result, err := q.NamedExec(query, c)
	if err != nil {
		return fmt.Errorf("failed to update contract: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// ContractExists checks if contract exists
func (q *ContractQueries) ContractExists(id uint) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM contracts WHERE contract_id = $1)`
	err := q.Get(&exists, query, id)
	if err != nil {
		return false, fmt.Errorf("failed to check contract existence: %w", err)
	}
	return exists, nil
}

// DeleteContract deletes contract by ID
func (q *ContractQueries) DeleteContract(id uint) error {
	query := `DELETE FROM contracts WHERE contract_id = $1`
	result, err := q.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete contract: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
