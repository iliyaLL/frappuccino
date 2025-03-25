package handlers

import (
	"errors"
	"frappuccino/internal/models"
	"frappuccino/internal/service"
	"frappuccino/internal/utils"
	"log/slog"
	"net/http"
)

type application struct {
	logger       *slog.Logger
	InventorySvc service.InventoryService
	// add more services
}

func NewApplication(inventorySvc service.InventoryService) *application {
	return &application{
		InventorySvc: inventorySvc,
		// add more services
	}
}

func mapErrorToResponse(err error, validationMap any) (int, any) {
	switch {
	// General errors
	case errors.Is(err, models.ErrInvalidID):
		return http.StatusBadRequest, utils.Response{"error": err.Error()}
	case errors.Is(err, models.ErrNoRecord):
		return http.StatusNotFound, utils.Response{"error": err.Error()}
	case errors.Is(err, models.ErrMissingFields):
		return http.StatusBadRequest, validationMap

	// Inventory errors
	case errors.Is(err, models.ErrDuplicateInventory),
		errors.Is(err, models.ErrNegativeQuantity),
		errors.Is(err, models.ErrInvalidEnumTypeInventory):
		return http.StatusBadRequest, utils.Response{"error": err.Error()}

	// Default catch-all
	default:
		return http.StatusInternalServerError, utils.Response{"error": "Internal Server Error"}
	}
}

func (app *application) Routes() http.Handler {
	router := http.NewServeMux()
	commonMiddleware := []Middleware{
		app.recoverPanic,
		app.logRequest,
	}

	endpoints := map[string]http.HandlerFunc{
		"POST /inventory":        app.inventoryCreate,
		"GET /inventory":         app.inventoryRetreiveAll,
		"GET /inventory/{id}":    app.inventoryRetrieveByID,
		"PUT /inventory/{id}":    app.inventoryUpdateByID,
		"DELETE /inventory/{id}": app.inventoryDeleteByID,
	}
	for endpoint, f := range endpoints {
		router.HandleFunc(endpoint, ChainMiddleware(f, commonMiddleware...))
	}

	return router
}
