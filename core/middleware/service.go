package middleware

import (
	commonSecurity "github.com/Alter17Ego/secuware"
	"github.com/Alter17Ego/secuware/http/middlewares/jwt"
	"github.com/xans-me/authopia/core/configuration"
)

// AuthService is IAuthService interface implementation
type AuthService struct {
	keycloak configuration.KeyCloak
}

// KeycloakAuthMiddleware to wrap routes api call with keycloak security module
func (svc AuthService) KeycloakAuthMiddleware() jwt.Middleware {
	uri := svc.keycloak.BaseURLAuth + "/auth/realms/" + svc.keycloak.Realm + KeycloakJwksValidationURL
	JWKsValidationMiddleware := commonSecurity.JwksValidationMiddleware(commonSecurity.JwksValidationMiddlewareOpts().WithDefaultOpts().SetJwksUrl(uri).AddWhitelist("/*").Build())
	return JWKsValidationMiddleware
}

// NewService function to init service instance
func NewService(keycloak configuration.KeyCloak) *AuthService {
	return &AuthService{keycloak: keycloak}
}
