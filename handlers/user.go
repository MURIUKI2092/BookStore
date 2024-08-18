package handlers

import (
	"BookStore/database"
	"BookStore/helpers"
	"BookStore/models"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// here is the struct to store the user
type UserResponse struct {
	UUID        uuid.UUID `json:"uuid,omitempty"`
	FirstName   string    `json:"first_name,omitempty"`
	LastName    string    `json:"last_name,omitempty"`
	PhoneNumber string    `json:"phone_number,omitempty"`
	Email       string    `json:"email,omitempty"`
	Role        string    `json:"role,omitempty"`
}

// ResponseStruct struct
type ResponseStruct struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// This creates a single user to the database
func CreateSingleUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}
	// check if the user is present
	userExists, err := helpers.UserEmailExists(user.Email)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "An error occurred while checking on user")
		return
	}
	if userExists {
		helpers.RespondWithError(w, http.StatusConflict, "User already exists")
		return
	}
	// here check if the phone number is present
	userPhonePresent, err := helpers.PhoneNumberPresent(user.PhoneNumber)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "An error occurred while checking on user")
		return
	}
	if userPhonePresent {
		helpers.RespondWithError(w, http.StatusConflict, "User already exists")
		return
	}
	// go ahead and hash the password before saving it
	err = user.HashPassword(user.Password)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "An error occurred while hashing the password")
		return
	}

	// go ahead and save the user in the database
	savedUser := database.DB.Create(&user)
	if savedUser.Error != nil {
		// throw an exception
		helpers.RespondWithError(w, http.StatusBadRequest, "An error occurred while creating a new user")
		return

	}
	// here is the response of the user
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

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(toReturnData); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}

}

// here we grab a single user
func GetSingleUser(w http.ResponseWriter, r *http.Request) {
	// here extract the uuid from the url
	userUuid := r.URL.Query().Get("uuid")
	if userUuid == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, "Uuid is required!")
		return

	}
	// fetch the single user here
	user, err := helpers.GrabSingleUserWithUuid(userUuid)
	// here handle all the errors
	if err != nil {
		if err.Error() == "invalid UUID format" {
			helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		} else if err.Error() == "record not found" {
			helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		} else {
			helpers.RespondWithError(w, http.StatusInternalServerError, "Error retrieving user")
			return
		}

	}

	// go ahead and return the user
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

// This function is to update the user
func UpdateSingleUser(w http.ResponseWriter, r *http.Request) {

	userUuid := r.URL.Query().Get("uuid")
	if userUuid == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, "Uuid is required!")
		return

	}
	// fetch the single user here
	requiredUser, err := helpers.GrabSingleUserWithUuid(userUuid)
	// here handle all the errors
	if err != nil {
		if err.Error() == "invalid UUID format" {
			helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		} else if err.Error() == "record not found" {
			helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		} else {
			helpers.RespondWithError(w, http.StatusInternalServerError, "Error retrieving user")
			return
		}

	}
	var updateData models.User
	// now I need to update the user here
	// Parse the incoming JSON data
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	// check if the user with the provided email and phone number is present
	var userEmailPresent models.User
	result := database.DB.Where("email = ? AND uuid <> ?", updateData.Email, userUuid).First(&userEmailPresent)
	if result.Error == nil {
		helpers.RespondWithError(w, http.StatusConflict, "Email is already in use by another user")
		return
	}

	// here do similar to the mobile phone
	var userPhonePresent models.User
	phone := database.DB.Where("phone_number = ? AND uuid <> ?", updateData.PhoneNumber, userUuid).First(&userPhonePresent)
	if phone.Error == nil {
		helpers.RespondWithError(w, http.StatusConflict, "Phone number is already in use by another user")
		return
	}
	// now go ahead and update the item
	requiredUser.PhoneNumber = updateData.PhoneNumber
	requiredUser.Email = updateData.Email
	requiredUser.Role = updateData.Role
	requiredUser.LastName = updateData.LastName
	requiredUser.FirstName = updateData.FirstName

	if err := database.DB.Save(&requiredUser).Error; err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Error updating user details")
		return
	}

	// return response
	// Prepare the response
	response := UserResponse{
		UUID:        requiredUser.UUID,
		FirstName:   requiredUser.FirstName,
		LastName:    requiredUser.LastName,
		PhoneNumber: requiredUser.PhoneNumber,
		Email:       requiredUser.Email,
		Role:        requiredUser.Role,
	}

	toReturnData := ResponseStruct{
		Msg:  "User updated successfully",
		Data: response,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(toReturnData); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}

}

// this service updates the user password only
func UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	var req models.UpdatePasswordRequest

	// Decode the request body into the struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}
	fmt.Printf("here is the req %v", req.ConfirmPassword)
	// Check if passwords match
	if req.Password != req.ConfirmPassword {
		helpers.RespondWithError(w, http.StatusBadRequest, "Passwords do not match")
		return
	}
	// here extract the uuid from the url
	userUuid := r.URL.Query().Get("uuid")
	if userUuid == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, "Uuid is required!")
		return

	}
	// fetch the single user here
	userPresent, err := helpers.GrabSingleUserWithUuid(userUuid)
	// here handle all the errors
	if err != nil {
		if err.Error() == "invalid UUID format" {
			helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		} else if err.Error() == "record not found" {
			helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		} else {
			helpers.RespondWithError(w, http.StatusInternalServerError, "Error retrieving user")
			return
		}

	}
	// here we check if the password is similar to the one already hashed
	err = userPresent.CheckPassword(req.PreviousPassword)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Password mismatch!")
		return
	}

	// here hash the user password
	err = userPresent.HashPassword(req.Password)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "An error occurred while hashing the password")
		return
	}
	if err := database.DB.Save(&userPresent).Error; err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Error updating user details")
		return
	}

	response := UserResponse{
		UUID:        userPresent.UUID,
		FirstName:   userPresent.FirstName,
		LastName:    userPresent.LastName,
		PhoneNumber: userPresent.PhoneNumber,
		Email:       userPresent.Email,
		Role:        userPresent.Role,
	}

	toReturnData := ResponseStruct{
		Msg:  "Password updated successfully",
		Data: response,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(toReturnData); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}

}
