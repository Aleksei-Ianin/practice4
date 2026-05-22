package repository

import (
	"calculator/api/internal/domain"
)

type CalculationRepository interface {
	Create(calc *domain.Calculation) error
	GetAll() ([]domain.Calculation, error)
	GetByID(id int) (*domain.Calculation, error)
	Delete(id int) error
}