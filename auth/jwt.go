package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secretkey")

var ExpiredTime = time.Now().Add(1 * time.Hour).Unix()
var ExpiredRefreshTime = time.Now().Add(time.Hour * 24 * 7).Unix()

type JWTClaim struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

func GenerateJWT(email string, expiresAt int64) (tokenString string, err error) {
	claims := &JWTClaim{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)

	return
}

func RefreshToken(accessTokenString string) (Tokens, error) {
	// Parse the JWT token
	refreshToken, err := jwt.Parse(accessTokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, http.ErrAbortHandler
		}
		return []byte(accessTokenString), nil
	})
	if err != nil || refreshToken == nil {
		return Tokens{}, fmt.Errorf("invalid refresh token")
	}

	if claims, ok := refreshToken.Claims.(jwt.MapClaims); ok && refreshToken.Valid {
		email, ok := claims["email"].(string)
		if !ok {
			return Tokens{}, fmt.Errorf("invalid email in refresh token claims")
		}

		newAccessToken, err := GenerateJWT(email, ExpiredTime)
		if err != nil {
			return Tokens{}, err
		}

		newRefreshAccessToken, err := GenerateJWT(email, ExpiredRefreshTime)
		if err != nil {
			return Tokens{}, err
		}

		return Tokens{AccessToken: newAccessToken, RefreshToken: newRefreshAccessToken}, nil
	}

	return Tokens{}, fmt.Errorf("invalid refresh token")
}

func ValidateToken(signedToken string) (err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return
	}

	return
}
