package utils

import (
	"encoding/json"
	"net/http"
)

func Message(status bool, message string) map[string]interface{} {

	return map[string]interface{}{"status": status, "message": message}

}

func Responds(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	_ = json.NewEncoder(w).Encode(data)
}
