package users

import (
	"strings"

	"github.com/Alter17Ego/go-dto-validator"
)

// AuthDataRequest struct
type AuthDataRequest struct {
	Email    string `json:"email" field:"required"`
	Password string `json:"password" field:"required" min:"6"`
}

// FormValidate method for AuthDataRequest
func (request *AuthDataRequest) FormValidate() (err error) {
	err = validator.Validate(*request)
	return
}

// UserLoginRequest struct
type UserLoginRequest struct {
	Username string `json:"username" field:"required"`
	Password string `json:"password" field:"required" min:"6"`
}

// FormValidate method for UserLoginRequest
func (request *UserLoginRequest) FormValidate() (err error) {
	err = validator.Validate(*request)
	return
}

// UserRegisterRequest struct
type UserRegisterRequest struct {
	AuthDataRequest
	Name        string `json:"name,omitempty"`
	PhoneNumber string `json:"phoneNumber" field:"required"`
}

// FormValidate method for UserRegisterRequest
func (request *UserRegisterRequest) FormValidate() (err error) {
	err = validator.Validate(*request)
	return
}

// UserPasswordChangeRequest struct
type UserPasswordChangeRequest struct {
	Email       string `json:"email" field:"required"`
	UserID      string `json:"userId" field:"required"`
	OldPassword string `json:"old_password" field:"required" min:"6"`
	NewPassword string `json:"new_password" field:"required" min:"6"`
}

// FormValidate method for UserPasswordChangeRequest
func (request *UserPasswordChangeRequest) FormValidate() (err error) {
	err = validator.Validate(*request)
	return
}

// UserIdentityRequest struct
type UserIdentityRequest struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
}

// FormValidate method for UserIdentityRequest
func (request *UserIdentityRequest) FormValidate() (err error) {
	err = validator.Validate(*request)
	return
}

// ConvertPhoneNumberToIndonesianArea function
func ConvertPhoneNumberToIndonesianArea(phoneNumber string) (convertedNumber string, err error) {

	if phoneNumber[:2] == "08" {
		convertedNumber = AreaPhoneNumbersIndonesian + phoneNumber[1:]
	} else if phoneNumber[:3] == "628" {
		convertedNumber = AreaPhoneNumbersIndonesian + phoneNumber[2:]
	} else if phoneNumber[:3] == AreaPhoneNumbersIndonesian {
		convertedNumber = phoneNumber
	} else {
		err = ErrInvalidPhoneNumber
	}

	return
}

// SplitName func to get a firstName and lastName from full name
func (request *UserRegisterRequest) SplitName() (firstName string, lastName string) {
	splitName := strings.Split(request.Name, " ")
	if len(splitName) == 1 {
		firstName = splitName[len(splitName)-1]
		lastName = ""
	} else {
		firstName = strings.Join(splitName[:len(splitName)-1], " ")
		lastName = splitName[len(splitName)-1]
	}
	return firstName, lastName
}
