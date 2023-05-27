package request

import (
	"encoding/json"
	"net/http"
)

func DecodeBody(r *http.Request, payload any) error {
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		return err
	}

	return nil
}
