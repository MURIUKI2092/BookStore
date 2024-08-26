package helpers

import (
	"BookStore/database"
	"BookStore/models"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// this checks if the user is present using phone numbers
func PhoneNumberPresent(PhoneNumber string) (bool, error) {
	var user models.User
	// check if the user is present
	userPresent := database.DB.Where("phone_number=?", PhoneNumber).First(&user)

	if userPresent.Error != nil {
		// if no user was found
		if userPresent.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		// if any other error return the error and a false
		return false, userPresent.Error

	}
	return true, nil
}

// This function checks if the user's email exists and avoids duplication in the database
func UserEmailExists(email string) (bool, error) {
	var user models.User
	// checks if the user is present in the database
	userPresent := database.DB.Where("email=?", email).First(&user)
	if userPresent.Error != nil {
		// if no user was found
		if userPresent.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		// if any other error return the error and a false
		return false, userPresent.Error

	}
	return true, nil

}

// This function checks if the user's email exists and avoids duplication in the database
func UserEmailPresent(email string) (*models.User, error) {
	var user models.User
	// checks if the user is present in the database
	userPresent := database.DB.Where("email=?", email).First(&user)
	if userPresent.Error != nil {
		// if no user was found
		if userPresent.Error != nil {
		// here handle when no record was found
		if userPresent.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("record not found")
		}
		return nil, userPresent.Error

	}

	}
	return &user, nil

}
// this grabs a single user record using the uuid
func GrabSingleUserWithUuid(userUuid string) (*models.User, error) {
	// check if the uuid is a valid uuid
	uuidValid, err := uuid.Parse(userUuid)
	if err != nil {
		return nil, errors.New("invalid UUID format")
	}

	var user models.User
	userPresent := database.DB.Where("uuid=?", uuidValid).First(&user)
	if userPresent.Error != nil {
		// here handle when no record was found
		if userPresent.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("record not found")
		}
		return nil, userPresent.Error

	}
	return &user, nil

}
