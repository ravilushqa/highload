package lib

import (
	"context"
	"net/http"

	"github.com/go-chi/jwtauth"

	"github.com/dgrijalva/jwt-go"
)

type UserClaims struct {
	jwt.StandardClaims
	UID int `json:"uid"`
}

type Auth struct {
	token *jwtauth.JWTAuth
}

func NewAuth(secret string) *Auth {
	return &Auth{token: jwtauth.New("HS256", []byte(secret), nil)}

}

func (a *Auth) GetToken() *jwtauth.JWTAuth {
	return a.token
}

func (a *Auth) EncodeToken(uid int) (string, error) {
	_, tokenString, err := a.token.Encode(UserClaims{
		UID: uid,
	})
	return tokenString, err
}

func (a *Auth) IsAuth(r *http.Request) bool {
	t, claims, err := jwtauth.FromContext(r.Context())
	_, ok := claims["UID"]
	return err == nil && t != nil && ok
}

func (a *Auth) GetUsedIDFromCtx(ctx context.Context) (int, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return 0, err
	}
	return claims["UID"].(int), nil
}
