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

func IsAuth(r *http.Request) bool {
	t, claims, err := jwtauth.FromContext(r.Context())
	_, ok := claims["uid"]
	return err == nil && t != nil && ok
}

func GetUsedIDFromCtx(ctx context.Context) (int, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return 0, err
	}
	return int(claims["uid"].(float64)), nil
}
