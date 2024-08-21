package handlers

import (
	"BookStore/database"
	"BookStore/helpers"
	"BookStore/models"
	"encoding/json"
	"fmt"
	"net/http"
)

// creates a single bookstore
func CreateSingleBookStore(w http.ResponseWriter, r *http.Request) {
	var store models.Store
	var user models.User
	userUuid := r.Context().Value("uuid").(string)

	err := json.NewDecoder(r.Body).Decode(&store)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}
	fmt.Printf("here is the phone number %v", store)
	// go ahead and check if there exist a store created by that user with either the phone or the email
	bookStorePresent, err := helpers.BookStoreExists(store.PhoneNumber, store.Email, userUuid)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "An error occurred while checking on the book store")
		return
	}

	if bookStorePresent {
		helpers.RespondWithError(w, http.StatusConflict, "The book store is already present")
		return
	}
	// run this items in a transaction
	// start the transaction
	tx := database.DB.Begin()
	store.CreatedBy = userUuid
	// now go ahead and create a  new one here
	savedStore := tx.Create(&store)
	if savedStore.Error != nil {
		tx.Rollback()
		// here throw an exception
		helpers.RespondWithError(w, http.StatusBadRequest, "An error occurred while creating a book store")
		return
	}

	// go ahead and update the user to have the book store uuid
	fmt.Printf("here is the store uuid: %v\n", store.UUID)
	// this need to happen in a transaction
	user.LinkBookStore = store.UUID
	// Update the user with the bookstore UUID
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		helpers.RespondWithError(w, http.StatusBadRequest, "An error occurred while updating the user with bookstore UUID")
		return
	}

	// Commit the transaction if all operations succeed
	if err := tx.Commit().Error; err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "An error occurred while committing the transaction")
		return
	}

	// Return the created store as a response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(store)

}

// Gets a single bookstore
func GetSingleBookStore(w http.ResponseWriter, r *http.Request) {
	storeUuid := r.URL.Query().Get("uuid")
	if storeUuid == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, "Uuid is required!")
		return

	}

	storePresent, err := helpers.GrabSingleStore(storeUuid)

	if err != nil {
		if err.Error() == "invalid UUId format" {
			helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		} else if err.Error() == "record not found" {
			helpers.RespondWithError(w, http.StatusNotFound, err.Error())
			return
		} else {
			helpers.RespondWithError(w, http.StatusInternalServerError, "Error retrieving book store")
			return
		}

	}

	toReturnData := ResponseStruct{
		Msg:  "Success",
		Data: storePresent,
	}
	w.Header().Set("Content-Type", "application/json")  
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(toReturnData); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}

}

// updates a single bookstore
func UpdateSingleBookStore(w http.ResponseWriter, r *http.Request) {
	userUuid := r.Context().Value("uuid").(string)
	bookStoreUuid := r.URL.Query().Get("uuid")
	if bookStoreUuid == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, "Uuid is required")
		return
	}
	fmt.Printf("holla ++++++++++>>>>>>>>>")

	// grab the single bookstore here
	requiredBookStore, err := helpers.GrabSingleStore(bookStoreUuid)
	if err != nil {
		if err.Error() == "invalid UUId format" {
			helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		} else if err.Error() == "record not found" {
			helpers.RespondWithError(w, http.StatusNotFound, err.Error())
			return
		} else {
			helpers.RespondWithError(w, http.StatusInternalServerError, "Error retrieving book store")
			return
		}

	}

	var updatedBookStore models.Store
	if err := json.NewDecoder(r.Body).Decode(&updatedBookStore); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}
	// here we check if there is another one present
	bookStorePresent, err := helpers.GrabSingleBookStore(updatedBookStore.PhoneNumber, updatedBookStore.Email, userUuid)
	fmt.Printf("here is the book store %v", bookStorePresent)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "An error occurred while checking on the book store")
		return
	}
	fmt.Printf("here is the book store %v", bookStorePresent)
	if bookStorePresent.UUID != requiredBookStore.UUID {
		helpers.RespondWithError(w, http.StatusConflict, "The book store is already present")
		return
	}

	requiredBookStore.City = updatedBookStore.City
	requiredBookStore.ContactPerson = updatedBookStore.ContactPerson
	requiredBookStore.ContactPersonEmail = updatedBookStore.ContactPersonEmail
	requiredBookStore.Description = updatedBookStore.Description
	requiredBookStore.Email = updatedBookStore.Email
	requiredBookStore.Name = updatedBookStore.Name
	requiredBookStore.Location = updatedBookStore.Location
	requiredBookStore.PhoneNumber = updatedBookStore.PhoneNumber
	requiredBookStore.ContactPersonPhoneNumber = updatedBookStore.ContactPersonPhoneNumber

	if err := database.DB.Save(&requiredBookStore).Error; err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Error updating Store details")
		return
	}

	toReturnData := ResponseStruct{
		Msg:  "Success",
		Data: requiredBookStore,
	}
	w.Header().Set("Content-Type", "application/json")  
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(toReturnData); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}

}
