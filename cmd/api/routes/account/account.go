package account

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/KiskaLE/RustDeskServer/cmd/api/database"
	"github.com/KiskaLE/RustDeskServer/cmd/api/middleware"
	"github.com/KiskaLE/RustDeskServer/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/valkey-io/valkey-glide/go/api"
	"github.com/valkey-io/valkey-glide/go/api/options"
	"golang.org/x/crypto/bcrypt"

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

	// chceck if account does not exist
	var account database.Accounts
	err = us.db.First(&account, "email = ?", payload.Email).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.WriteError(w, http.StatusBadRequest, errors.New("account already exists"))
		return
	}

	// hash password
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, errors.New("Failed to hash password: "+err.Error()))
		return
	}
	// create account
	newaccount := &database.Accounts{
		Email:    payload.Email,
		Password: string(hashedPassword),
	}
	err = us.db.Save(newaccount).Error
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, errors.New("Failed to create account: "+err.Error()))
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "account created"})
}

func (us *AccountService) LoginRoute(w http.ResponseWriter, r *http.Request) {
	var payload LoginPayload
	err := utils.ParseJSON(r, &payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// check if account exists
	var account database.Accounts
	err = us.db.First(&account, "email = ?", payload.Email).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Email or password is incorrect"))
		return
	}

	// check if password is correct
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(payload.Password))
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Email or password is incorrect"))
		return
	}

	// generate tokens
	tokenString, refreshToken, err := generateTokens(strconv.Itoa(int(account.ID)))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// save refresh token to valkey
	err = us.saveRefreshToken(refreshToken, account)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": tokenString, "refresh_token": refreshToken})
	return
}

func (us *AccountService) LogoutRoute(w http.ResponseWriter, r *http.Request) {
	// TODO implement
	// invalidate refresh token
	_, err := us.valkey.Del([]string{"refresh_token:" + r.Header.Get("refresh_token")})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// blacklist jwt token
	claimsValue := r.Context().Value("claims")

	keyOptions := options.SetOptions{
		Expiry: &options.Expiry{
			Count: uint64((claimsValue.(*middleware.Claims).ExpiresAt.Unix()) - time.Now().Unix()),
			Type:  options.Seconds,
		},
	}
	_, err = us.valkey.SetWithOptions("jwt_blacklist:"+claimsValue.(*middleware.Claims).Jti, "true", keyOptions)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "logout successful"})
}

func (us *AccountService) RefreshTokenRoute(w http.ResponseWriter, r *http.Request) {
	// TODO implement
	var payload RefreshTokenPayload
	err := utils.ParseJSON(r, &payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// get refresh token from valkey
	oldRefreshToken, err := us.valkey.Get("refresh_token:" + payload.RefreshToken)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("invalid refresh token"))
		return
	}

	// invalidate refresh token
	_, err = us.valkey.Del([]string{"refresh_token:" + payload.RefreshToken})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// find account
	var account database.Accounts
	err = us.db.First(&account, "id = ?", oldRefreshToken.Value()).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("invalid refresh token"))
		return
	}

	// generate tokens
	tokenString, refreshToken, err := generateTokens(strconv.Itoa(int(account.ID)))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// save refresh token to valkey
	err = us.saveRefreshToken(refreshToken, account)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": tokenString, "refresh_token": refreshToken})
	return
}

func generateRefreshToken() (string, error) {
	b := make([]byte, 64)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (us *AccountService) saveRefreshToken(refreshToken string, account database.Accounts) error {
	opts := options.SetOptions{
		Expiry: &options.Expiry{
			Count: 60 * 60 * 24, // 24 hours
			Type:  options.Seconds,
		},
	}

	// save refresh token to valkey
	_, err := us.valkey.SetWithOptions("refresh_token:"+refreshToken, strconv.Itoa(int(account.ID)), opts)
	if err != nil {
		return err
	}

	return nil
}

func generateTokens(id string) (tokenString string, refreshToken string, err error) {
	// Generate JWT token
	expirationTime := time.Now().Add(30 * time.Minute)
	claims := &middleware.Claims{
		Jti: uuid.New().String(),
		ID:  id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtKey := []byte(os.Getenv("SALT"))

	tokenString, err = token.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	// generate refresh token
	refreshToken, err = generateRefreshToken()
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshToken, nil
}
