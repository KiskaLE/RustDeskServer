package computer

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/KiskaLE/RustDeskServer/cmd/api/db"
	"github.com/KiskaLE/RustDeskServer/utils"
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
	err = conn.First(&computer, "name = ?", strings.ToLower(payload.ComputerName)).Error
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
	}

	// Update data about computer
	computer.RustDeskID = payload.RustDeskID
	computer.IP = payload.IP
	computer.OS = payload.OS
	computer.OSVersion = payload.OSVersion
	computer.LastConnection = sql.NullTime{Time: time.Now(), Valid: true}

	err = conn.Save(computer).Error
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Computer refreshed"})
	return
}
