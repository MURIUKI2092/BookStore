package handlers

import (
	"BookStore/database"
	"BookStore/helpers"
	"BookStore/models"
	"encoding/json"
	"fmt"
	"net/http"
)

// This creates a single book
func CreatesSingleBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	userUuid := r.Context().Value("uuid").(string)
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}
	userPresent, err := helpers.GrabSingleUserWithUuid(userUuid)
	fmt.Printf("user present: %v\n", userPresent)

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
	bookStoreUuid := userPresent.LinkBookStore
	fmt.Printf("userpresent %v", userPresent)
	// Convert UUID to string
	bookStoreUuidStr := bookStoreUuid.String()

	// go ahead and grab the book
	bookPresent, err := helpers.IsBookPresent(book.Title, bookStoreUuidStr)
	if err != nil {
		// Handle the error appropriately
		helpers.RespondWithError(w, http.StatusInternalServerError, "Error checking if the book is present")
		return
	}

	if bookPresent {
		// Handle the case where the book is already present
		helpers.RespondWithError(w, http.StatusConflict, "The book is already present in the store")
		return
	}

	// go ahead and add the book
	book.CreatedBy = userUuid
	book.RemainingQuantity = book.Quantity
	book.Store = bookStoreUuidStr
	savedBook := database.DB.Create(&book)
	if savedBook.Error != nil {
		fmt.Printf("here is the error: %v", savedBook.Error)
		helpers.RespondWithError(w, http.StatusBadRequest, "An error occurred while creating a single book")
		return
	}

	toReturnData := ResponseStruct{
		Msg:  "Success",
		Data: book,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(toReturnData); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}

}

// Gets all the books in a book store
func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	userUuid := r.Context().Value("uuid").(string)

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

	userStore := user.LinkBookStore.String()

	//grab all the books here
	allBooks, err := helpers.AllStoreBooks(userStore)

	if err != nil {
		if err.Error() == "record not found" {
			helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		} else {
			helpers.RespondWithError(w, http.StatusInternalServerError, "Error retrieving user")
			return
		}
	}

	toReturnData := ResponseStruct{
		Msg:  "Success",
		Data: allBooks,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(toReturnData); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}

}

// Get a single book from a book store
func GetSingleBook(w http.ResponseWriter, r *http.Request) {
	userUuid := r.Context().Value("uuid").(string)
	// fetch the single user here
	bookUuid := r.URL.Query().Get("uuid")
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

	// user store
	userStore := user.LinkBookStore.String()
	bookPresent, err := helpers.SingleBook(userStore, bookUuid)
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

	toReturnData := ResponseStruct{
		Msg:  "Success",
		Data: bookPresent,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(toReturnData); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}

}

// updates a single book
func UpdateSingleBook(w http.ResponseWriter, r *http.Request) {
	var updateData models.Book
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}
	userUuid := r.Context().Value("uuid").(string)
	// fetch the single user here
	bookUuid := r.URL.Query().Get("uuid")
	user, err := helpers.GrabSingleUserWithUuid(userUuid)
	// here handle all the errors
	fmt.Printf("here %v", user)
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

	// user store
	userStore := user.LinkBookStore.String()
	book, err := helpers.SingleBook(userStore, bookUuid)
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
	// go ahead and update the book
	
	book.Author = updateData.Author
	book.Edition = updateData.Edition
	book.Genre = updateData.Genre
	book.ISBN = updateData.ISBN
	book.Title = updateData.Title
	book.Quantity = updateData.Quantity
	book.RemainingQuantity = updateData.RemainingQuantity
	book.Pages = updateData.Pages
	book.Langage = updateData.Langage
	book.Publisher = updateData.Publisher
	book.PublicationDate = updateData.PublicationDate

	if err := database.DB.Save(&book).Error; err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Error updating user details")
		return
	}

	toReturnData := ResponseStruct{
		Msg:  "Book Updated successfully",
		Data: book,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(toReturnData); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}
}
