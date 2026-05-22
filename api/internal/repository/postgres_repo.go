package repository

import (
	"calculator/api/internal/domain"
	"database/sql"
	"errors"
	"log"
)

type PostgresCalculationRepository struct {
	db *sql.DB
}

func NewPostgresCalculationRepository(db *sql.DB) *PostgresCalculationRepository {
	return &PostgresCalculationRepository{db: db}
}

func (r *PostgresCalculationRepository) Create(calc *domain.Calculation) error {
	query := `INSERT INTO calculations (operand1, operand2, operation, result, created_at) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRow(query, calc.Operand1, calc.Operand2, calc.Operation, calc.Result, calc.CreatedAt).Scan(&calc.ID)
	if err != nil {
		log.Printf("Create error: %v", err)
		return err
	}
	return nil
}

func (r *PostgresCalculationRepository) GetAll() ([]domain.Calculation, error) {
	rows, err := r.db.Query(`SELECT id, operand1, operand2, operation, result, created_at 
	                         FROM calculations ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var calcs []domain.Calculation
	for rows.Next() {
		var c domain.Calculation
		err := rows.Scan(&c.ID, &c.Operand1, &c.Operand2, &c.Operation, &c.Result, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		calcs = append(calcs, c)
	}
	return calcs, nil
}

func (r *PostgresCalculationRepository) GetByID(id int) (*domain.Calculation, error) {
	var c domain.Calculation
	query := `SELECT id, operand1, operand2, operation, result, created_at 
	          FROM calculations WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&c.ID, &c.Operand1, &c.Operand2, &c.Operation, &c.Result, &c.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *PostgresCalculationRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM calculations WHERE id = $1`, id)
	return err
}