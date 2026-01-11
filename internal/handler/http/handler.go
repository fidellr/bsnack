package http

import (
	"bsnack/internal/domain"
	"bsnack/internal/service"
	"encoding/json"
	"net/http"
)

type Handler struct {
	prodSvc  *service.ProductService
	transSvc *service.TransactionService
	custSvc  *service.CustomerService
}

func NewHandler(prodSvc *service.ProductService, transSvc *service.TransactionService, custSvc *service.CustomerService) *Handler {
	return &Handler{
		prodSvc:  prodSvc,
		transSvc: transSvc,
		custSvc:  custSvc,
	}
}

func (h *Handler) respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}

func (h *Handler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]string{"error": message})
}

// Customer Handlers

// GET /customers
func (h *Handler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	customers, err := h.custSvc.GetAllCustomers(r.Context())
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.respondJSON(w, http.StatusOK, customers)
}

// Product Handlers

// POST /products
func (h *Handler) AddProduct(w http.ResponseWriter, r *http.Request) {
	var p domain.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.prodSvc.AddProduct(r.Context(), &p); err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusCreated, p)
}

// GET /products?date=2025-10-22
func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		h.respondError(w, http.StatusBadRequest, "date query param required (YYYY-MM-DD)")
		return
	}

	products, err := h.prodSvc.GetProductsByDate(r.Context(), date)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, products)
}

// Transaction Handlers

// POST /transactions
func (h *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req service.PurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.transSvc.Purchase(r.Context(), req); err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusCreated, map[string]string{"message": "Transaction successful"})
}

// POST /redemptions (Redeem Points)
func (h *Handler) Redeem(w http.ResponseWriter, r *http.Request) {
	type RedeemReq struct {
		CustomerName string `json:"customer_name"`
		ProductID    int64  `json:"product_id"`
	}
	var req RedeemReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.transSvc.Redeem(r.Context(), req.CustomerName, req.ProductID); err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Redemption successful"})
}

// GET /transactions?start=2025-10-01&end=2025-12-31
func (h *Handler) GetReport(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	if start == "" {
		h.respondError(w, http.StatusBadRequest, "start date required")
		return
	}

	report, err := h.transSvc.GetReport(r.Context(), start, end)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, report)
}
