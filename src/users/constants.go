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
	ErrGeneral            = response.ErrorStruct{ErrorCode: 1, Message: "internal server error"}
	ErrDecodeRequestData  = response.ErrorStruct{ErrorCode: 2, Message: "decode request data failed"}
	ErrBusinessLogic      = response.ErrorStruct{ErrorCode: 3, Message: "business logic execution failed"}
	ErrUserCredential     = response.ErrorStruct{ErrorCode: 4, Message: "credential mismatch"}
	ErrEmailUnverified    = response.ErrorStruct{ErrorCode: 5, Message: "email is not yet verified"}
	ErrGetUserInfo        = response.ErrorStruct{ErrorCode: 6, Message: "get user info failed"}
	ErrVerifyOTP          = response.ErrorStruct{ErrorCode: 7, Message: "otp code is invalid"}
	ErrUniqueUser         = response.ErrorStruct{ErrorCode: 8, Message: "users must be unique"}
	ErrActionEmail        = response.ErrorStruct{ErrorCode: 9, Message: "execute action email failed"}
	ErrWrongPassword      = response.ErrorStruct{ErrorCode: 10, Message: "wrong password"}
	ErrPhoneIsExist       = response.ErrorStruct{ErrorCode: 11, Message: "phone number is already exist"}
	ErrEmailIsExist       = response.ErrorStruct{ErrorCode: 12, Message: "email is already exist"}
	ErrCantSendOTP        = response.ErrorStruct{ErrorCode: 13, Message: "cannot resend otp yet"}
	ErrExpiredOTP         = response.ErrorStruct{ErrorCode: 14, Message: "otp code is expired"}
	ErrInvalidPhoneNumber = response.ErrorStruct{ErrorCode: 15, Message: "invalid phone number"}
	ErrSamePhoneNumber    = response.ErrorStruct{ErrorCode: 16, Message: "cannot use old phone numbers"}
	ErrCouchDBNotFound    = response.ErrorStruct{ErrorCode: 17, Message: "DB : user not found / available to creating new"}
	ErrExistCouchDBUser   = response.ErrorStruct{ErrorCode: 18, Message: "DB : user already exist"}
)
