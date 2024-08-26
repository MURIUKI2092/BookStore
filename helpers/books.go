package helpers

import (
	"BookStore/database"
	"BookStore/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func BookPresent(Title string, Store string) (*models.Book, error) {
	var book models.Book

	bookPresent := database.DB.Where("store =? AND title=?", Store, Title).First(&book)

	if bookPresent.Error != nil {
		// handle when no record was present
		if bookPresent.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("record not found")
		}
		return nil, bookPresent.Error
	}
	return &book, nil
}

func IsBookPresent(Title string, Store string) (bool, error) {
	var book models.Book

	bookPresent := database.DB.Where("store =? AND title=?", Store, Title).First(&book)

	if bookPresent.Error != nil {
		// if no book present
		if bookPresent.Error == gorm.ErrRecordNotFound {

			return false, nil
		}
		return false, bookPresent.Error
	}
	return true, nil

}
func AllStoreBooks(storeUuid string) ([]models.Book, error) {
	var store models.Store
	var books []models.Book
	storePresent := database.DB.Where("uuid", storeUuid).First(&store)
	if storePresent.Error != nil {
		if storePresent.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("record not found")
		}
		return nil, storePresent.Error
	}

	// grab all the books here
	result := database.DB.Where("store = ?", storeUuid).Find(&books)
	if result.Error != nil {
		return nil, result.Error
	}

	return books, nil

}
func SingleBook(storeUuid string, bookUuid string) (*models.Book, error) {
	uuidValid, err := uuid.Parse(bookUuid)
	if err != nil {
		return nil, errors.New("invalid UUID format")
	}
	var book models.Book

	bookPresent := database.DB.Where("uuid=? AND store =?", uuidValid, storeUuid).First(&book)
	if bookPresent.Error != nil {
		// here handle when no record was found
		if bookPresent.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("record not found")
		}
		return nil, bookPresent.Error

	}
	return &book, nil
}
