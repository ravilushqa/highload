package lib

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/jwtauth"

	"github.com/dgrijalva/jwt-go"
)

type UserClaims struct {
	jwt.StandardClaims
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
	_, tokenString, err := a.token.Encode(jwt.StandardClaims{
		Subject:   strconv.Itoa(uid),
		ExpiresAt: time.Now().AddDate(0, 0, 14).Unix(),
	})
	return tokenString, err
}

func IsAuth(r *http.Request) bool {
	t, claims, err := jwtauth.FromContext(r.Context())
	_, ok := claims["sub"]
	return err == nil && t != nil && ok
}

func GetAuthUserID(ctx context.Context) (int, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(claims["sub"].(string))
}
