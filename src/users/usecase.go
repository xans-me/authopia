package users

import (
	"context"
	"errors"
	"github.com/Nerzal/gocloak/v7"
	"github.com/xans-me/authopia/core/configuration"
	"strings"
)

// UserService struct
type UserService struct {
	repository IUsersRepository
	keycloak   configuration.KeyCloak
	//monitoring      monitoring.ErrorMonitoring
}

// Login service function
func (svc UserService) Login(context context.Context, request UserLoginRequest) (data *gocloak.JWT, err error) {

	// data validation
	err = request.FormValidate()
	if err != nil {
		return
	}

	// converted phone input with area number
	if request.Username[:2] == "08" || request.Username[:3] == "628" || request.Username[:3] == AreaPhoneNumbersIndonesian {
		request.Username, _ = ConvertPhoneNumberToIndonesianArea(request.Username)
	}

	return svc.repository.LoginUserKeycloak(context, request)
}

// Register service function
func (svc UserService) Register(context context.Context, request UserRegisterRequest) (data *gocloak.JWT, err error) {

	// request data validation
	err = request.FormValidate()
	if err != nil {
		return
	}

	request.PhoneNumber, err = ConvertPhoneNumberToIndonesianArea(request.PhoneNumber)
	if err != nil {
		return
	}

	// checking if data is already exist
	err = svc.repository.CheckingIsExistUsers(context, UserIdentityRequest{
		Email:       request.Email,
		PhoneNumber: request.PhoneNumber,
	})
	if err != nil {
		return
	}

	// user registration to keycloak
	_, err = svc.repository.RegisterUserKeycloak(context, request)
	if err != nil {
		result := strings.SplitAfterN(err.Error(), " ", 3)
		// send message of error
		err = errors.New(result[2])
		return
	}

	//// send email verify if user registration successfully
	//_, err = svc.repository.ExecuteResendVerifyEmail(context, request.Email)
	//if err != nil {
	//	result := strings.SplitAfterN(err.Error(), " ", 3)
	//	// send message of error
	//	return nil, errors.New(result[2])
	//}

	// build data for user auto login after register
	loginCredential := UserLoginRequest{
		ClientID:     svc.keycloak.ClientID,
		ClientSecret: svc.keycloak.ClientSecret,
		Username:     request.Email,
		Password:     request.Password,
	}

	// do auto login to response a access token
	return svc.repository.LoginUserKeycloak(context, loginCredential)
}

// NewService function to init service instance
func NewService(repository IUsersRepository, keycloak configuration.KeyCloak) *UserService {
	return &UserService{repository: repository, keycloak: keycloak}
}
