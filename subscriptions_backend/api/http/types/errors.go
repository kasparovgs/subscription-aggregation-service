package types

import (
	"encoding/json"
	"net/http"
	"subscriptions_backend/domain"
)

func ProcessError(w http.ResponseWriter, err error, resp any) {
	if err != nil {
		if myErr, ok := err.(*domain.MyErr); ok {
			http.Error(w, myErr.Error(), myErr.Code)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	if resp != nil {
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
		}
	}
}
