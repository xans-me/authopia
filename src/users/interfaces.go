//go:generate mockgen -source ./interfaces.go -package users -destination ./mock_user.go

package users

import (
	"context"
	"github.com/Nerzal/gocloak/v7"
)

// IUseCase service interfaces
type IUseCase interface {
	Login(ctx context.Context, request UserLoginRequest) (*gocloak.JWT, error)
	Register(context context.Context, request UserRegisterRequest) (*gocloak.JWT, error)
}

// IUsersRepository repository interfaces
type IUsersRepository interface {
	LoginAdminKeycloak(ctx context.Context, request UserLoginRequest) (interface{}, error)
	LoginUserKeycloak(ctx context.Context, request UserLoginRequest) (*gocloak.JWT, error)
	RegisterUserKeycloak(ctx context.Context, request UserRegisterRequest) (interface{}, error)
	ChangePasswordUserKeycloak(ctx context.Context, request UserPasswordChangeRequest) error
	UpdateUserKeycloak(ctx context.Context, userData gocloak.User) error
	GetCredentialUserKeycloak(ctx context.Context, UUID string) ([]*gocloak.CredentialRepresentation, error)
	ExecuteForgotPassword(ctx context.Context, email string) (interface{}, error)
	ExecuteResendVerifyEmail(ctx context.Context, email string) (interface{}, error)
	AssignUserToGroupKeycloak(ctx context.Context, userID string) error
	CheckingIsExistUsers(ctx context.Context, request UserIdentityRequest) error
}
