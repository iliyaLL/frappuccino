package repository

import (
	"database/sql"
	"errors"
	"frappuccino/internal/models"
	"log/slog"

	"github.com/lib/pq"
)

type InventoryRepository interface {
	Insert(name, unit string, quantity int, categories []string) error
	RetrieveByID(id int) (models.Inventory, error)
	RetrieveAll() ([]models.Inventory, error)
	Update(id int, name, unit string, quantity int, categories []string) error
	Delete(id int) error
}

type inventoryRepositoryPostgres struct {
	pq     *sql.DB
	logger *slog.Logger
}

func NewInventoryRepositoryWithPostgres(db *sql.DB, logger *slog.Logger) *inventoryRepositoryPostgres {
	return &inventoryRepositoryPostgres{
		pq:     db,
		logger: logger,
	}
}

func (m *inventoryRepositoryPostgres) Insert(name, unit string, quantity int, categories []string) error {
	_, err := m.pq.Exec("INSERT INTO inventory (name, quantity, unit, categories) VALUES ($1, $2, $3, $4)", quantity, unit, pq.Array(categories))
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				return models.ErrDuplicateInventory
			case "23514":
				return models.ErrNegativeQuantity
			case "22P02":
				return models.ErrInvalidEnumTypeInventory
			}
		}
		return err
	}

	return nil
}

func (m *inventoryRepositoryPostgres) RetrieveByID(id int) (models.Inventory, error) {
	var inventory models.Inventory
	err := m.pq.QueryRow("SELECT * FROM inventory WHERE id = $1", id).Scan(
		&inventory.ID,
		&inventory.Name,
		&inventory.Quantity,
		&inventory.Unit,
		pq.Array(&inventory.Categories),
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Inventory{}, models.ErrNoRecord
		}
		return models.Inventory{}, err
	}

	return inventory, nil
}

func (m *inventoryRepositoryPostgres) RetrieveAll() ([]models.Inventory, error) {
	rows, err := m.pq.Query("SELECT * FROM inventory")
	if err != nil {
		m.logger.Error("Failed to execute Query", "error", err)
		return nil, err
	}
	defer rows.Close()

	var InventoryAll []models.Inventory
	for rows.Next() {
		var inventory models.Inventory

		err = rows.Scan(
			&inventory.ID,
			&inventory.Name,
			&inventory.Quantity,
			&inventory.Unit,
			pq.Array(&inventory.Categories),
		)
		if err != nil {
			return nil, err
		}

		InventoryAll = append(InventoryAll, inventory)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return InventoryAll, err
}

func (m *inventoryRepositoryPostgres) Update(id int, name, unit string, quantity int, categories []string) error {
	result, err := m.pq.Exec("UPDATE inventory SET name=$1, unit=$2, quantity=$3, categories=$4 WHERE id=$5", name, unit, quantity, pq.Array(categories), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrNoRecord
		}
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				return models.ErrDuplicateInventory
			case "23514":
				return models.ErrNegativeQuantity
			case "22P02":
				return models.ErrInvalidEnumTypeInventory
			}
		}

		return err
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return models.ErrNoRecord
	}

	return err
}

func (m *inventoryRepositoryPostgres) Delete(id int) error {
	result, err := m.pq.Exec("DELETE FROM inventory WHERE id=$1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return models.ErrNoRecord
	}

	return err
}
