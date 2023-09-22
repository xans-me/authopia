//go:generate mockgen -source ./interfaces.go -package users -destination ./mock_user.go

package users

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
)

// IUseCase service interfaces
type IUseCase interface {
	Login(ctx context.Context, request UserLoginRequest) (*gocloak.JWT, error)
	Register(context context.Context, request UserRegisterRequest) (*gocloak.JWT, error)
}

// IUsersRepository repository interfaces
type IUsersRepository interface {
	LoginUserKeycloak(ctx context.Context, request UserLoginRequest) (data *gocloak.JWT, err error)
	RegisterUserKeycloak(ctx context.Context, request UserRegisterRequest) (data interface{}, err error)
	ChangePasswordUserKeycloak(ctx context.Context, request UserPasswordChangeRequest) (err error)
	UpdateUserKeycloak(ctx context.Context, userData gocloak.User) (err error)
	GetCredentialUserKeycloak(ctx context.Context, UUID string) (representations []*gocloak.CredentialRepresentation, err error)
	ExecuteForgotPassword(ctx context.Context, email string) (data interface{}, err error)
	ExecuteResendVerifyEmail(ctx context.Context, email string) (data interface{}, err error)
	AssignUserToGroupKeycloak(ctx context.Context, userID string) (err error)
	CheckingIsExistUsers(ctx context.Context, request UserIdentityRequest) (err error)
}
