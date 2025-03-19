package service

import (
	"database/sql"
	"frappuccino/internal/models"
	"frappuccino/internal/repository"
	"log/slog"
)

type InventoryService interface {
	Insert(name, unit string, quantity int, categories *[]string) error
	RetrieveByID(id int32) (models.Inventory, error)
	RetrieveAll() (*[]models.Inventory, error)
}

type inventoryService struct {
	inventoryRepo repository.InventoryRepository
}

func NewInventoryService(db *sql.DB, logger *slog.Logger) *inventoryService {
	return &inventoryService{
		inventoryRepo: repository.NewInventoryRepositoryWithPostgres(db, logger),
	}
}

func (s *inventoryService) Insert(name, unit string, quantity int, categories *[]string) error {
	return nil
}
func (s *inventoryService) RetrieveByID(id int32) (models.Inventory, error) {
	return models.Inventory{}, nil
}
func (s *inventoryService) RetrieveAll() (*[]models.Inventory, error) {
	return nil, nil
}
