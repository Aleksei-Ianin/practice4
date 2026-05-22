package service

import (
	"calculator/api/internal/domain"
	"calculator/api/internal/repository"
	"errors"
	"math"
	"time"
)

type CalculationService struct {
	repo repository.CalculationRepository
}

func NewCalculationService(repo repository.CalculationRepository) *CalculationService {
	return &CalculationService{repo: repo}
}

// Perform вычисляет результат и сохраняет запись
func (s *CalculationService) Perform(operand1, operand2 float64, operation string) (*domain.Calculation, error) {
	var result float64
	switch operation {
	case "add":
		result = operand1 + operand2
	case "subtract":
		result = operand1 - operand2
	case "multiply":
		result = operand1 * operand2
	case "divide":
		if operand2 == 0 {
			return nil, errors.New("division by zero")
		}
		result = operand1 / operand2
	case "power":
		result = math.Pow(operand1, operand2)
	default:
		return nil, errors.New("unknown operation")
	}

	calc := &domain.Calculation{
		Operand1:  operand1,
		Operand2:  operand2,
		Operation: operation,
		Result:    result,
		CreatedAt: time.Now(),
	}
	err := s.repo.Create(calc)
	if err != nil {
		return nil, err
	}
	return calc, nil
}

func (s *CalculationService) GetAll() ([]domain.Calculation, error) {
	return s.repo.GetAll()
}

func (s *CalculationService) Delete(id int) error {
	return s.repo.Delete(id)
}