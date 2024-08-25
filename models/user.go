package models

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName     string    `gorm:"size:100;not null"`
	LastName      string    `gorm:"size:100;not null"`
	PhoneNumber   string    `gorm:"size:15;unique;not null"`
	Email         string    `gorm:"size:100;unique;not null"`
	Password      string    `gorm:"not null"`
	Role          string    `gorm:"not null"`
	UUID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	LinkBookStore uuid.UUID `gorm:"type:uuid;"`
}

type UpdatePasswordRequest struct {
	PreviousPassword string `json:"previous_password"`
	Password         string `json:"password"`
	ConfirmPassword  string `json:"confirm_password"`
}

type Claims struct {
	UUID string `json:"uuid"`
	jwt.StandardClaims
}

// create a function to hash the user's password as you are saving the user
func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

// create a function to check if the password is a match
// CheckPassword checks if the provided password matches the hashed password
func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}

// generates uuid before creation
func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.UUID = uuid.New()
	return
}
