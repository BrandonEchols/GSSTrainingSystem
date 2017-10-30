package controllers

import (
	"encoding/json"
	"net/http"
)

type controller struct{}

func (this controller) WriteErrorMessageWithStatus(w http.ResponseWriter, code int, msg string) {
	payload := map[string]interface{}{
		"error": msg,
	}

	marshaled_data, _ := json.Marshal(payload)
	w.Write(marshaled_data)
	w.WriteHeader(code)
}
