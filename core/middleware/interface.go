package middleware

import (
	"github.com/Alter17Ego/secuware/http/middlewares/jwt"
)

// IAuthService interface
type IAuthService interface {
	KeycloakAuthMiddleware() jwt.Middleware
}
