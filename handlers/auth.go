package handlers

import (
	"BookStore/helpers"
	"BookStore/models"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("your_secret_key")

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginResponse struct {
	Token string `json:"token"`
}

// This is a function that handles the login for a new user
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var req LoginUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Input")
		return
	}

	// grab the user with that email
	userExists, err := helpers.UserEmailPresent(req.Email)
	if err != nil {
		if err.Error() == "record not found" {
			helpers.RespondWithError(w, http.StatusNotFound, err.Error())
			return
		} else {

			helpers.RespondWithError(w, http.StatusBadRequest, "An error occurred while checking on user")
			return
		}
	}

	// go ahead and check if the user password matches the hashed password
	err = userExists.CheckPassword(req.Password)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Incorrect credentials!")
		return
	}

	// here go ahead and perform a login
	// create a new jwt
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &models.Claims{
		UUID: userExists.UUID.String(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	// go ahead and sign the jwt key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	// Return the token in the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LoginResponse{Token: tokenString})

}

func GrabUserProfile(w http.ResponseWriter, r *http.Request) {
	userUuid := r.Context().Value("uuid").(string)
	user, err := helpers.GrabSingleUserWithUuid(userUuid)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "User not found")
		return
	}

	response := UserResponse{
		UUID:        user.UUID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Role:        user.Role,
	}

	toReturnData := ResponseStruct{
		Msg:  "Success",
		Data: response,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(toReturnData); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}
}
