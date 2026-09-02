package products

import (
	"log"
	"net/http"
	"strconv"

	repo "github/M-b-a-s/myStore/internal/adapters/postgresql/sqlc"
	"github/M-b-a-s/myStore/internal/json"

	"github.com/go-chi/chi/v5"
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

func (h *handler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	// 1. Extract the product ID from the request (assuming it's passed as a query parameter)
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		json.Write(w, http.StatusBadRequest, nil, "Missing product ID")
		return
	}

	// Convert the ID to int64
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		json.Write(w, http.StatusBadRequest, nil, "Invalid product ID")
		return
	}

	// 2. Call the service -> GetProductByID
	product, err := h.service.GetProductByID(r.Context(), id)
	if err != nil {
		log.Printf("Error fetching product: %v", err)
		json.Write(w, http.StatusInternalServerError, nil, "Error fetching product")
		return
	}

	json.Write(w, http.StatusOK, product, "Product fetched successfully")
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
