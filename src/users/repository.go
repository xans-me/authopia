package users

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/xans-me/authopia/core/configuration"
)

// Repository struct
type Repository struct {
	keycloak configuration.KeyCloak
}

// LoginAdminKeycloak repo func
func (repo Repository) LoginAdminKeycloak(ctx context.Context, request UserLoginRequest) (data interface{}, err error) {
	client := gocloak.NewClient(repo.keycloak.BaseURLAuth)
	data, err = client.LoginAdmin(ctx, request.Username, request.Password, repo.keycloak.Realm)
	return
}

// LoginUserKeycloak repo func
func (repo Repository) LoginUserKeycloak(ctx context.Context, request UserLoginRequest) (data *gocloak.JWT, err error) {
	client := gocloak.NewClient(repo.keycloak.BaseURLAuth)
	return client.Login(ctx, request.ClientID, request.ClientSecret, repo.keycloak.Realm, request.Username, request.Password)
}

// RegisterUserKeycloak repo func
func (repo Repository) RegisterUserKeycloak(ctx context.Context, request UserRegisterRequest) (data interface{}, err error) {
	client := gocloak.NewClient(repo.keycloak.BaseURLAuth)
	token, err := client.LoginAdmin(ctx, repo.keycloak.AdminUsername, repo.keycloak.AdminPassword, repo.keycloak.Realm)
	if err != nil {
		return nil, err
	}

	// split a name into firstName and lastName
	firstName, lastName := request.SplitName()

	user := gocloak.User{
		FirstName:       &firstName,
		LastName:        &lastName,
		Email:           &request.Email,
		Enabled:         &DefaultUserEnabledKeycloak,
		EmailVerified:   &DefaultUserEmailVerified,
		Username:        &request.PhoneNumber,
		RequiredActions: &[]string{},
	}

	// register user to keycloak
	data, err = client.CreateUser(ctx, token.AccessToken, repo.keycloak.Realm, user)
	if err != nil {
		return nil, err
	}

	// set user password
	err = client.SetPassword(ctx, token.AccessToken, data.(string), repo.keycloak.Realm, request.Password, false)
	return
}

// ChangePasswordUserKeycloak repo func
func (repo Repository) ChangePasswordUserKeycloak(ctx context.Context, request UserPasswordChangeRequest) (err error) {
	client := gocloak.NewClient(repo.keycloak.BaseURLAuth)
	token, err := client.LoginAdmin(ctx, repo.keycloak.AdminUsername, repo.keycloak.AdminPassword, repo.keycloak.Realm)
	if err != nil {
		return err
	}

	// set user password
	err = client.SetPassword(ctx, token.AccessToken, request.UserID, repo.keycloak.Realm, request.NewPassword, false)
	return
}

// UpdateUserKeycloak repo func
func (repo Repository) UpdateUserKeycloak(context context.Context, userData gocloak.User) (err error) {
	client := gocloak.NewClient(repo.keycloak.BaseURLAuth)
	token, err := client.LoginAdmin(context, repo.keycloak.AdminUsername, repo.keycloak.AdminPassword, repo.keycloak.Realm)
	if err != nil {
		return err
	}
	err = client.UpdateUser(context, token.AccessToken, repo.keycloak.Realm, userData)
	return
}

// GetCredentialUserKeycloak repo func
func (repo Repository) GetCredentialUserKeycloak(context context.Context, UUID string) (representations []*gocloak.CredentialRepresentation, err error) {
	client := gocloak.NewClient(repo.keycloak.BaseURLAuth)
	token, err := client.LoginAdmin(context, repo.keycloak.AdminUsername, repo.keycloak.AdminPassword, repo.keycloak.Realm)
	if err != nil {
		return representations, err
	}

	return client.GetCredentials(context, token.AccessToken, repo.keycloak.Realm, UUID)
}

// ExecuteForgotPassword to execute action of forgotten password
func (repo Repository) ExecuteForgotPassword(ctx context.Context, email string) (data interface{}, err error) {
	client := gocloak.NewClient(repo.keycloak.BaseURLAuth)
	token, err := client.LoginAdmin(ctx, repo.keycloak.AdminUsername, repo.keycloak.AdminPassword, repo.keycloak.Realm)

	if err != nil {
		return nil, err
	}

	users, err := client.GetUsers(ctx, token.AccessToken, repo.keycloak.Realm, gocloak.GetUsersParams{Email: &email})
	if err != nil {
		return nil, ErrGetUserInfo
	}

	// getting fists users
	if len(users) == 0 || len(users) > 1 {
		return nil, ErrUniqueUser
	}
	user := users[0]

	err = client.ExecuteActionsEmail(ctx, token.AccessToken, repo.keycloak.Realm, gocloak.ExecuteActionsEmail{
		UserID:   user.ID,
		Lifespan: &LifespanActionEmail,
		Actions:  &[]string{ActionExecuteUpdatePassword},
	})
	if err != nil {
		return nil, ErrActionEmail
	}

	data = user.Email
	return
}

// ExecuteResendVerifyEmail to execute action resend verify email of users
func (repo Repository) ExecuteResendVerifyEmail(ctx context.Context, email string) (data interface{}, err error) {
	client := gocloak.NewClient(repo.keycloak.BaseURLAuth)
	token, err := client.LoginAdmin(ctx, repo.keycloak.AdminUsername, repo.keycloak.AdminPassword, repo.keycloak.Realm)
	if err != nil {
		return nil, err
	}

	users, err := client.GetUsers(ctx, token.AccessToken, repo.keycloak.Realm, gocloak.GetUsersParams{Email: &email})
	if err != nil {
		return nil, ErrGetUserInfo
	}

	// getting fists users
	if len(users) == 0 || len(users) > 1 {
		return nil, ErrUniqueUser
	}
	user := users[0]

	err = client.ExecuteActionsEmail(ctx, token.AccessToken, repo.keycloak.Realm, gocloak.ExecuteActionsEmail{
		UserID:   user.ID,
		Lifespan: &LifespanActionEmail,
		Actions:  &[]string{ActionExecuteVerifyEmail},
	})
	if err != nil {
		return nil, ErrActionEmail
	}

	data = *user.Email
	return
}

// AssignUserToGroupKeycloak to assign user to groups in keycloak
func (repo Repository) AssignUserToGroupKeycloak(ctx context.Context, userID string) (err error) {
	// Assign User to Groups => Senyumku lite users
	client := gocloak.NewClient(repo.keycloak.BaseURLAuth)
	token, err := client.LoginAdmin(ctx, repo.keycloak.AdminUsername, repo.keycloak.AdminPassword, repo.keycloak.Realm)
	if err != nil {
		return
	}

	dataGroup, err := client.GetGroups(ctx, token.AccessToken, repo.keycloak.Realm, gocloak.GetGroupsParams{Search: &SenyumkuLiteGroupKeycloak})
	if err != nil {
		return
	}

	var groupID string
	for _, valueGroup := range dataGroup {
		if strings.EqualFold(*valueGroup.Name, SenyumkuLiteGroupKeycloak) {
			groupID = *valueGroup.ID
		} else {
			break
		}
	}

	err = client.AddUserToGroup(ctx, token.AccessToken, repo.keycloak.Realm, userID, groupID)
	if err != nil {
		return
	}

	// Add Custom Attributes => [ senyumku-lite.verified , senyumku-lite.verifiedAt ]
	err = client.UpdateUser(ctx, token.AccessToken, repo.keycloak.Realm, gocloak.User{
		ID:      &userID,
		Enabled: &DefaultUserEnabledKeycloak,
		Attributes: &map[string][]string{
			"senyumku-lite.verified":   {"true"},
			"senyumku-lite.verifiedAt": {strconv.Itoa(int(time.Now().UnixNano() / int64(time.Millisecond)))},
		},
	})

	return
}

// CheckingIsExistUsers is a repo func to get user info from keycloak
func (repo Repository) CheckingIsExistUsers(ctx context.Context, request UserIdentityRequest) (err error) {
	client := gocloak.NewClient(repo.keycloak.BaseURLAuth)
	dataJWT, err := client.LoginAdmin(ctx, repo.keycloak.AdminUsername, repo.keycloak.AdminPassword, repo.keycloak.Realm)
	if err != nil {
		return
	}

	// checking if phone number is already used
	var userByPhone []*gocloak.User

	if request.PhoneNumber != "" {
		userByPhone, err = client.GetUsers(ctx, dataJWT.AccessToken, repo.keycloak.Realm, gocloak.GetUsersParams{
			Username: &request.PhoneNumber,
		})
		if err != nil {
			return
		}

		if len(userByPhone) > 0 {
			return ErrPhoneIsExist
		}
	}

	// checking if email is already used
	var userByEmail []*gocloak.User
	if request.Email != "" {
		userByEmail, err = client.GetUsers(ctx, dataJWT.AccessToken, repo.keycloak.Realm, gocloak.GetUsersParams{
			Email: &request.Email,
		})
		if err != nil {
			return
		}

		if len(userByEmail) > 0 {
			return ErrEmailIsExist
		}
	}

	return
}

// NewRepository to create new repository instance
func NewRepository(keycloak configuration.KeyCloak) *Repository {
	return &Repository{keycloak: keycloak}
}
