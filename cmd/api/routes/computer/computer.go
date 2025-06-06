package computer

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/KiskaLE/RustDeskServer/cmd/api/db"
	"github.com/KiskaLE/RustDeskServer/utils"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RefreshComputerRoute(w http.ResponseWriter, r *http.Request) {
	var payload RefreshComputerPayload
	err := utils.ParseJSON(r, &payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	conn := db.Connect()

	// Check if computer exist
	computer := &db.Computers{}
	err = conn.First(&computer, "rust_desk_id = ?", payload.RustDeskID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create computer
		newComputer := &db.Computers{
			Name:           strings.ToLower(payload.ComputerName),
			RustDeskID:     payload.RustDeskID,
			IP:             payload.IP,
			OS:             payload.OS,
			OSVersion:      payload.OSVersion,
			LastConnection: sql.NullTime{Time: time.Now(), Valid: true},
		}
		err = conn.Save(newComputer).Error
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Computer created"})
		return
	}

	// Update data about computer
	err = conn.Model(&computer).Updates(db.Computers{
		Name:           strings.ToLower(payload.ComputerName),
		IP:             payload.IP,
		OS:             payload.OS,
		OSVersion:      payload.OSVersion,
		LastConnection: sql.NullTime{Time: time.Now(), Valid: true},
	}).Error

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Computer refreshed"})
	return
}

func GetComputerRustDeskIDRoute(w http.ResponseWriter, r *http.Request) {
	computerName := mux.Vars(r)["computerName"]

	conn := db.Connect()

	computer := &db.Computers{}
	err := conn.First(&computer, "name = ?", computerName).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	} else if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, computer.RustDeskID)
	return
}
