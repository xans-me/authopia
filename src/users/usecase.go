package users

import (
	"context"
	"errors"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/xans-me/authopia/core/configuration"
)

// UseCase struct
type UseCase struct {
	repository IUsersRepository
	keycloak   configuration.KeyCloak
	//monitoring      monitoring.ErrorMonitoring
}

// Login service function
func (svc UseCase) Login(context context.Context, request UserLoginRequest) (data *gocloak.JWT, err error) {

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
func (svc UseCase) Register(context context.Context, request UserRegisterRequest) (data *gocloak.JWT, err error) {

	// request data validation
	err = request.FormValidate()
	if err != nil {
		return
	}

	request.PhoneNumber, err = ConvertPhoneNumberToIndonesianArea(request.PhoneNumber)
	if err != nil {
		return
	}

	// checking if data is already existing
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
		// send a message of error
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

// NewUseCase function to init useCase instance
func NewUseCase(repository IUsersRepository, keycloak configuration.KeyCloak) *UseCase {
	return &UseCase{repository: repository, keycloak: keycloak}
}
