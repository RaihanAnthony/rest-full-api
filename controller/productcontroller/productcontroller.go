package productcontroller 

import (
	"net/http"

	"rest_full_api/helper"
)

func Index(w http.ResponseWriter, r *http.Request) {

	data := []map[string]interface{}{
		{
			"id": 1,
			"nama": "kemeja",
			"stock": 100,
		},
		{
			"id": 2,
			"nama": "polo",
			"stock": 200,
		},
		{
			"id": 3,
			"nama": "celana",
			"stock": 250,
		},
	}

	helper.ResponseJSON(w, http.StatusOK, data)
} 