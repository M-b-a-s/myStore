package products

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	repo "github/M-b-a-s/myStore/internal/adapters/postgresql/sqlc"
	"github/M-b-a-s/myStore/internal/json"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
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
	id, err := parseProductID(r)
	if err != nil {
		json.Write(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	// 2. Call the service -> GetProductByID
	product, err := h.service.GetProductByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			json.Write(w, http.StatusNotFound, nil, "Product not found")
			return
		}
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

func (h *handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := parseProductID(r)
	if err != nil {
		json.Write(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	// 2. Parse the request body into a Product struct
	var product repo.UpdateProductParams
	if err := json.Read(r, &product); err != nil {
		log.Printf("Error reading request body: %v", err)
		json.Write(w, http.StatusBadRequest, err, "Invalid request body")
		return
	}
	product.ID = id

	// 3. Call the service -> UpdateProduct
	updatedProduct, err := h.service.UpdateProduct(r.Context(), product)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			json.Write(w, http.StatusNotFound, nil, "Product not found")
			return
		}
		log.Printf("Error updating product: %v", err)
		json.Write(w, http.StatusInternalServerError, err, "Error updating product")
		return
	}

	json.Write(w, http.StatusOK, updatedProduct, "Product updated successfully")
}

func (h *handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// 1. Extract the product ID from the request
	id, err := parseProductID(r)
	if err != nil {
		json.Write(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	// 2. Call the service -> DeleteProduct
	_, err = h.service.DeleteProduct(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			json.Write(w, http.StatusNotFound, nil, "Product not found")
			return
		}
		log.Printf("Error deleting product: %v", err)
		json.Write(w, http.StatusInternalServerError, err, "Error deleting product")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseProductID(r *http.Request) (int64, error) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		return 0, errors.New("missing product ID")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid product ID")
	}

	return id, nil
}
