package helpers

import (
	"BookStore/database"
	"BookStore/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func BookStoreExists(PhoneNumber string, Email string, UserUuid string) (bool, error) {
	var store models.Store

	storePresent := database.DB.Where("created_by = ? OR (phone_number = ? AND email = ?)",
		UserUuid, PhoneNumber, Email).First(&store)

	if storePresent.Error != nil {
		//  if no store was present
		if storePresent.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		// if any other error return the error and a false
		return false, storePresent.Error
	}
	return true, nil

}

func GrabSingleBookStore(PhoneNumber string, Email string, UserUuid string) (*models.Store, error) {
	var store models.Store

	storePresent := database.DB.Where("created_by = ? OR (phone_number = ? AND email = ?)",
		UserUuid, PhoneNumber, Email).First(&store)

	if storePresent.Error != nil {
		// handle when no record is present
		if storePresent.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("record not found")
		}
		return nil, storePresent.Error

	}
	return &store, nil

}

func GrabSingleStore(userUuid string) (*models.Store, error) {
	uuidValid, err := uuid.Parse(userUuid)
	if err != nil {
		return nil, errors.New("Invalid UUID format")
	}

	var store models.Store
	storePresent := database.DB.Where("uuid=? ", uuidValid).First(&store)
	if storePresent.Error != nil {
		// handle when no record is present
		if storePresent.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("record not found")
		}
		return nil, storePresent.Error

	}
	return &store, nil
}
