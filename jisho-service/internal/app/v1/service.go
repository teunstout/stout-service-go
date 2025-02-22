package v1

import (
	"encoding/json"
	"net/http"
)

func HandleJishoReponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	kw := r.URL.Query().Get("keyword")
	if kw == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jr, err := http.Get("https://jisho.org/api/v1/search/words?keyword=" + kw)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	var response JishoResponse
	err = json.NewDecoder(jr.Body).Decode(&response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
