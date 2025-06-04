package test

import (
	"net/http"

	"github.com/KiskaLE/RustDeskServer/utils"
)

func HelloRoute(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Hello World!"})
	return
}
