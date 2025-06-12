package account

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/KiskaLE/RustDeskServer/cmd/api/database"
	"github.com/KiskaLE/RustDeskServer/utils"
	"github.com/valkey-io/valkey-glide/go/api"

	"gorm.io/gorm"
)

type AccountService struct {
	db     *gorm.DB
	valkey api.GlideClientCommands
}

func NewAccountService(db *gorm.DB, valkey api.GlideClientCommands) *AccountService {
	return &AccountService{db: db, valkey: valkey}
}

func (us *AccountService) RegisterRoute(w http.ResponseWriter, r *http.Request) {
	// TODO implement
	var payload RegisterPayload
	err := utils.ParseJSON(r, &payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// chceck if user does not exist
	var user database.Accounts
	err = us.db.First(&user, "email = ?", payload.Email).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.WriteError(w, http.StatusBadRequest, errors.New("User already exists"))
		return
	}

	// hash password
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, errors.New("Failed to hash password: "+err.Error()))
		return
	}
	// create user
	newUser := &database.Accounts{
		Email:    payload.Email,
		Password: string(hashedPassword),
	}
	err = us.db.Save(newUser).Error
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, errors.New("Failed to create user: "+err.Error()))
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "User created"})
}

func (us *AccountService) LoginRoute(w http.ResponseWriter, r *http.Request) {
	// TODO implement
	// SET key
	_, err := us.valkey.Set("key", "val")
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, errors.New("Failed to set key: "+err.Error()))
		return
	}

	val, err := us.valkey.Get("key")
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, errors.New("Failed to get key: "+err.Error()))
		return
	}

	fmt.Println(val.Value())

	utils.WriteJSON(w, http.StatusOK, val.Value())
	return
}
