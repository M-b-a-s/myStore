package products

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	repo "github/M-b-a-s/myStore/internal/adapters/postgresql/sqlc"
	"github/M-b-a-s/myStore/internal/json"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type handler struct {
	service Service
}

func validateProductInput(name string, priceInCents, quantity int32) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("product name is required")
	}
	if priceInCents < 0 {
		return errors.New("price_in_cents must not be negative")
	}
	if quantity < 0 {
		return errors.New("quantity must not be negative")
	}

	return nil
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
	var input createProductRequest
	if err := json.Read(r, &input); err != nil {
		if errors.Is(err, io.EOF) {
			json.Write(w, http.StatusBadRequest, nil, "Request body is required")
			return
		}
		log.Printf("Error reading request body: %v", err)
		json.Write(w, http.StatusBadRequest, nil, "Invalid request body")
		return
	}
	if err := validateProductInput(input.Name, input.PriceInCents, input.Quantity); err != nil {
		json.Write(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	// 2. Call the service -> CreateProduct
	createdProduct, err := h.service.CreateProduct(r.Context(), repo.CreateProductParams{
		Name:         input.Name,
		PriceInCents: input.PriceInCents,
		Quantity:     input.Quantity,
	})
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

	var input updateProductRequest
	if err := json.Read(r, &input); err != nil {
		if errors.Is(err, io.EOF) {
			json.Write(w, http.StatusBadRequest, nil, "Request body is required")
			return
		}
		log.Printf("Error reading request body: %v", err)
		json.Write(w, http.StatusBadRequest, nil, "Invalid request body")
		return
	}
	if err := validateProductInput(input.Name, input.PriceInCents, input.Quantity); err != nil {
		json.Write(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	// 3. Call the service -> UpdateProduct
	updatedProduct, err := h.service.UpdateProduct(r.Context(), repo.UpdateProductParams{
		ID:           id,
		Name:         input.Name,
		PriceInCents: input.PriceInCents,
		Quantity:     input.Quantity,
	})
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
