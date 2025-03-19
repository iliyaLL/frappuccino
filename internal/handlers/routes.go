package handlers

import (
	"frappuccino/internal/service"
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

func (app *application) Routes() http.Handler {
	router := http.NewServeMux()
	commonMiddleware := []Middleware{
		app.recoverPanic,
		app.logRequest,
		contentTypeJSON,
	}

	endpoints := map[string]http.HandlerFunc{
		"POST /inventory": app.inventoryCreatePost,
	}
	for endpoint, f := range endpoints {
		router.HandleFunc(endpoint, ChainMiddleware(f, commonMiddleware...))
	}

	return router
}
