package handlers

import (
	"frappuccino/internal/service"
	"log/slog"
	"net/http"
)

type application struct {
	logger       *slog.Logger
	InventorySvc service.InventoryService
	MenuSvc      service.MenuService
	// add more services
}

func NewApplication(inventorySvc service.InventoryService, menuSvc service.MenuService) *application {
	return &application{
		InventorySvc: inventorySvc,
		MenuSvc:      menuSvc,
		// add more services
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

		"POST /menu":        app.menuCreate,
		"GET /menu":         app.menuRetrieveAll,
		"GET /menu/{id}":    app.menuRetrieveAllByID,
		"PUT /menu/{id}":    app.menuUpdate,
		"DELETE /menu/{id}": app.menuDelete,
	}
	for endpoint, f := range endpoints {
		router.HandleFunc(endpoint, ChainMiddleware(f, commonMiddleware...))
	}

	return router
}
