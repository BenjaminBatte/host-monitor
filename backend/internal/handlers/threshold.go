package threshold

import (
	"encoding/json"
	"net/http"

	"github.com/BenjaminBatte/host-monitor/internal/config"
)

type ThresholdRequest struct {
	Threshold int `json:"threshold"`
}

func ThresholdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var req ThresholdRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		if req.Threshold < 1 || req.Threshold > 10000 {
			http.Error(w, "Invalid threshold range", http.StatusBadRequest)
			return
		}

		if err := config.UpdateSettings(config.Settings{LatencyThreshold: req.Threshold}); err != nil {
			http.Error(w, "Failed to update threshold", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodGet {
		current := config.GetThreshold()
		json.NewEncoder(w).Encode(map[string]int{"threshold": current})
		return
	}

	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}
