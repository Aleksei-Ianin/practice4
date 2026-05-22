package repository

import (
	"calculator/api/internal/domain"
	"testing"
	"time"
)

// In-memory реализация для тестов
type mockRepo struct {
	calcs map[int]domain.Calculation
	nextID int
}

func (m *mockRepo) Create(c *domain.Calculation) error {
	c.ID = m.nextID
	m.nextID++
	m.calcs[c.ID] = *c
	return nil
}
func (m *mockRepo) GetAll() ([]domain.Calculation, error) {
	list := []domain.Calculation{}
	for _, v := range m.calcs {
		list = append(list, v)
	}
	return list, nil
}
func (m *mockRepo) GetByID(id int) (*domain.Calculation, error) {
	c, ok := m.calcs[id]
	if !ok {
		return nil, nil
	}
	return &c, nil
}
func (m *mockRepo) Delete(id int) error {
	delete(m.calcs, id)
	return nil
}

func TestMockRepoCreate(t *testing.T) {
	repo := &mockRepo{calcs: make(map[int]domain.Calculation)}
	calc := &domain.Calculation{Operand1: 2, Operand2: 3, Operation: "add", Result: 5, CreatedAt: time.Now()}
	err := repo.Create(calc)
	if err != nil {
		t.Fatal(err)
	}
	if calc.ID != 0 {
		t.Log("ID assigned correctly")
	}
	if len(repo.calcs) != 1 {
		t.Errorf("Expected 1 calculation, got %d", len(repo.calcs))
	}
}

func TestMockRepoGetAll(t *testing.T) {
	repo := &mockRepo{calcs: make(map[int]domain.Calculation)}
	
	calc1 := &domain.Calculation{Operand1: 1, Operand2: 2, Operation: "add", Result: 3}
	if err := repo.Create(calc1); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	
	calc2 := &domain.Calculation{Operand1: 4, Operand2: 5, Operation: "multiply", Result: 20}
	if err := repo.Create(calc2); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	
	all, err := repo.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 2 {
		t.Errorf("Expected 2, got %d", len(all))
	}
}

func TestMockRepoDelete(t *testing.T) {
	repo := &mockRepo{calcs: make(map[int]domain.Calculation)}
	calc := &domain.Calculation{Operand1: 7, Operand2: 8, Operation: "subtract", Result: -1}
	if err := repo.Create(calc); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if err := repo.Delete(0); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	all, _ := repo.GetAll()
	if len(all) != 0 {
		t.Error("Delete failed")
	}
}