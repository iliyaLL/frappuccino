package models

import "errors"

var (
	ErrNoRecord         = errors.New("models: no record")
	ErrNegativeQuantity = errors.New("models: positive quantity constraint violation")
	ErrNegativePrice    = errors.New("models: positive price constraint violation")
	ErrMissingFields    = errors.New("models: missing fields")
	ErrInvalidID        = errors.New("id is not valid int")

	ErrDuplicateInventory       = errors.New("models: duplicate inventory")
	ErrInvalidEnumTypeInventory = errors.New("models: invalid enum type. Supported types: shots, ml, g, units")
)
