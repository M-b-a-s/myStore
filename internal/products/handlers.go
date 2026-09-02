package products

import (
	"log"
	"net/http"

	repo "github/M-b-a-s/myStore/internal/adapters/postgresql/sqlc"
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
		json.Write(w, http.StatusInternalServerError, err, "Error listing products")
		return
	}

	json.Write(w, http.StatusOK, products, "Products listed successfully")
}

func (h *handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	// 1. Parse the request body into a Product struct
	var product repo.CreateProductParams
	if err := json.Read(r, &product); err != nil {
		log.Printf("Error reading request body: %v", err)
		json.Write(w, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	// 2. Call the service -> CreateProduct
	createdProduct, err := h.service.CreateProduct(r.Context(), product)
	if err != nil {
		log.Printf("Error creating product: %v", err)
		json.Write(w, http.StatusInternalServerError, err, "Error creating product")
		return
	}

	json.Write(w, http.StatusCreated, createdProduct, "Product created successfully")
}
