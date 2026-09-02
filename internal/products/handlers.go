package products

import (
	"log"
	"net/http"

	"github/M-b-a-s/myStore/internal/json"
)

type handler struct {
	service Service
}

func NewHandler(s Service) *handler {
	return &handler{
		service: s,
	}
}

func (h *handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	// 1. Call the service -> ListProducts
	products, err := h.service.ListProducts(r.Context())
	if err != nil {
		log.Printf("Error listing products: %v", err)
		json.Write(w, http.StatusInternalServerError, err)
		return
	}

	json.Write(w, http.StatusOK, products)
}
