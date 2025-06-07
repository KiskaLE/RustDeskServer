package computer

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/KiskaLE/RustDeskServer/cmd/api/database"
	"github.com/KiskaLE/RustDeskServer/utils"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type ComputerService struct {
	db *gorm.DB
}

func NewComputerService(db *gorm.DB) *ComputerService {
	return &ComputerService{db: db}
}

func (cs *ComputerService) RefreshComputerRoute(w http.ResponseWriter, r *http.Request) {
	var payload RefreshComputerPayload
	err := utils.ParseJSON(r, &payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Check if computer exist
	computer := &database.Computers{}
	err = cs.db.First(&computer, "rust_desk_id = ?", payload.RustDeskID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create computer
		newComputer := &database.Computers{
			Name:           strings.ToLower(payload.ComputerName),
			RustDeskID:     payload.RustDeskID,
			IP:             payload.IP,
			OS:             payload.OS,
			OSVersion:      payload.OSVersion,
			LastConnection: sql.NullTime{Time: time.Now(), Valid: true},
		}
		err = cs.db.Save(newComputer).Error
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Computer created"})
		return
	}

	// Update data about computer
	err = cs.db.Model(&computer).Updates(database.Computers{
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

func (cs *ComputerService) GetComputerRustDeskIDRoute(w http.ResponseWriter, r *http.Request) {
	computerName := mux.Vars(r)["computerName"]

	computer := &database.Computers{}
	err := cs.db.First(&computer, "name = ?", computerName).Error
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
