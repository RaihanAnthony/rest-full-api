package helper 

import (
	"encoding/json"
	"net/http"
)

func ResponseJSON(w http.ResponseWriter, code int, playload interface{}) {
	response, _ := json.Marshal(playload)
	w.Header().Add("content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}