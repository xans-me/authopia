package users

import "github.com/xans-me/authopia/helpers/response"

// user domain constants
const (
	MessageOK                   string = "ok"
	KeycloakJwksValidationURL   string = "/protocol/openid-connect/certs"
	CouchDBUserPrefix           string = "usr"
	ActionExecuteVerifyEmail    string = "VERIFY_EMAIL"
	ActionExecuteUpdatePassword string = "UPDATE_PASSWORD"
	AreaPhoneNumbersIndonesian  string = "+62"
)

var (
	AuthopiaGroupKeycloak      = "Authopia users"
	AuthopiaRolesSigner        = "ROLE_SIGNER"
	LifespanActionEmail        = 500
	DefaultUserEnabledKeycloak = true
	DefaultUserEmailVerified   = false
)

// error variable default values
var (
	ErrGeneral            = response.ErrorStruct{Code: 1, Message: "internal server error"}
	ErrDecodeRequestData  = response.ErrorStruct{Code: 2, Message: "decode request data failed"}
	ErrBusinessLogic      = response.ErrorStruct{Code: 3, Message: "business logic execution failed"}
	ErrUserCredential     = response.ErrorStruct{Code: 4, Message: "credential mismatch"}
	ErrEmailUnverified    = response.ErrorStruct{Code: 5, Message: "email is not yet verified"}
	ErrGetUserInfo        = response.ErrorStruct{Code: 6, Message: "get user info failed"}
	ErrVerifyOTP          = response.ErrorStruct{Code: 7, Message: "otp code is invalid"}
	ErrUniqueUser         = response.ErrorStruct{Code: 8, Message: "users must be unique"}
	ErrActionEmail        = response.ErrorStruct{Code: 9, Message: "execute action email failed"}
	ErrWrongPassword      = response.ErrorStruct{Code: 10, Message: "wrong password"}
	ErrPhoneIsExist       = response.ErrorStruct{Code: 11, Message: "phone number is already exist"}
	ErrEmailIsExist       = response.ErrorStruct{Code: 12, Message: "email is already exist"}
	ErrCantSendOTP        = response.ErrorStruct{Code: 13, Message: "cannot resend otp yet"}
	ErrExpiredOTP         = response.ErrorStruct{Code: 14, Message: "otp code is expired"}
	ErrInvalidPhoneNumber = response.ErrorStruct{Code: 15, Message: "invalid phone number"}
	ErrSamePhoneNumber    = response.ErrorStruct{Code: 16, Message: "cannot use old phone numbers"}
	ErrCouchDBNotFound    = response.ErrorStruct{Code: 17, Message: "DB : user not found / available to creating new"}
	ErrExistCouchDBUser   = response.ErrorStruct{Code: 18, Message: "DB : user already exist"}
)
