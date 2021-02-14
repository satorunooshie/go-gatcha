package auth

import (
	_"github.com/dgrijalva/jwt-go"
	jwt "github.com/form3tech-oss/jwt-go"
	middleware "github.com/auth0/go-jwt-middleware"
	"net/http"
	"time"
)

var GetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// header
	token := jwt.New(jwt.SigningMethodHS256)
	// calims
	claims := token.Claims.(jwt.MapClaims)
	claims["admin"] = true
	claims["sub"] = "54546557354"
	claims["name"] = "username"
	claims["iat"] = time.Now()
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// digital signature
	tokenString, _ := token.SignedString([]byte("secret"))
	//return jwt
	w.Write([]byte(tokenString))
})

// JwtMiddleware check token
var JwtMiddleware = middleware.New(middleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})