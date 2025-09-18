package service

import (
	"bxs/log"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

// PriceResponse represents the response structure for price API endpoints
type PriceResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ClosestPriceData represents the data structure for closest price response
type ClosestPriceData struct {
	Price       string `json:"price"`
	Timestamp   int64  `json:"timestamp"`
	Tolerance   int64  `json:"tolerance"`
	RequestTime int64  `json:"request_time"`
}

// LatestPriceData represents the data structure for latest price response
type LatestPriceData struct {
	Price       string `json:"price"`
	Timestamp   int64  `json:"timestamp"`
	RequestTime int64  `json:"request_time"`
}

func (ps *priceService) StartApiServer(port int) {
	// Create a new mux for routing
	mux := http.NewServeMux()

	// Register API endpoints
	mux.HandleFunc("/api/price/closest", ps.handleGetClosestPrice)
	mux.HandleFunc("/api/price/latest", ps.handleGetLatestPrice)

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PriceResponse{
			Success: true,
			Data:    map[string]string{"status": "ok"},
		})
	})

	// Default handler for root path
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PriceResponse{
			Success: true,
			Data: map[string]interface{}{
				"service": "price-service",
				"endpoints": map[string]string{
					"closest_price": "/api/price/closest?timestamp=<timestamp>&tolerance=<tolerance>",
					"latest_price":  "/api/price/latest",
					"health":        "/health",
				},
			},
		})
	})

	go func() {
		log.Logger.Info("Price service HTTP server starting", zap.Int("port", port))
		if err := http.ListenAndServe(fmt.Sprintf("%s:%d", "0.0.0.0", port), mux); err != nil {
			log.Logger.Fatal("HTTP server failed to start", zap.Error(err))
		}
	}()
}

// handleGetClosestPrice handles the GET request for closest price by timestamp
func (ps *priceService) handleGetClosestPrice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse query parameters
	timestampStr := r.URL.Query().Get("timestamp")
	toleranceStr := r.URL.Query().Get("tolerance")

	if timestampStr == "" {
		json.NewEncoder(w).Encode(PriceResponse{
			Success: false,
			Error:   "timestamp parameter is required",
		})
		return
	}

	// Parse timestamp
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		json.NewEncoder(w).Encode(PriceResponse{
			Success: false,
			Error:   "invalid timestamp format",
		})
		return
	}

	// Parse tolerance (optional, default to 300 seconds)
	tolerance := int64(300)
	if toleranceStr != "" {
		tolerance, err = strconv.ParseInt(toleranceStr, 10, 64)
		if err != nil {
			json.NewEncoder(w).Encode(PriceResponse{
				Success: false,
				Error:   "invalid tolerance format",
			})
			return
		}
	}

	// Get closest price
	price, err := ps.GetClosestPriceByTimestamp(timestamp, tolerance)
	if err != nil {
		json.NewEncoder(w).Encode(PriceResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Return success response
	json.NewEncoder(w).Encode(PriceResponse{
		Success: true,
		Data: ClosestPriceData{
			Price:       price.String(),
			Timestamp:   timestamp,
			Tolerance:   tolerance,
			RequestTime: time.Now().Unix(),
		},
	})
}

// handleGetLatestPrice handles the GET request for latest price
func (ps *priceService) handleGetLatestPrice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get latest price
	price, timestamp, err := ps.GetLatestPrice()
	if err != nil {
		json.NewEncoder(w).Encode(PriceResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Return success response
	json.NewEncoder(w).Encode(PriceResponse{
		Success: true,
		Data: LatestPriceData{
			Price:       price.String(),
			Timestamp:   timestamp,
			RequestTime: time.Now().Unix(),
		},
	})
}
