package types

import "github.com/golang-jwt/jwt/v5"

// UserAPI representa un usuario en el JWT.
type UserClaims struct {
	UserAPI
	jwt.RegisteredClaims
}
