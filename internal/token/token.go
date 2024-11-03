package token

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Token struct {
	Login string
	Role  string
}

func Generate(login, role, secret string, accessTime time.Duration) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"login": login,
			"role":  role,
			"exp":   time.Now().Add(accessTime).Unix(),
		})
	accessToken, err := t.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func Verify(tokenStringBearer, secret string) (jwt.MapClaims, error) {
	splitToken := strings.Split(tokenStringBearer, " ")
	if len(splitToken) != 2 {
		return nil, fmt.Errorf("token is invalid (not bearer)")
	}
	tokenString := splitToken[1]

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	// expTimeInterface, ok := claims["exp"]
	// if !ok {
	// 	return fmt.Errorf("token is invalid (can`t parse [exp] claims)")
	// }
	// expTimeUNIX, ok := expTimeInterface.(int64)
	// if !ok {
	// 	return fmt.Errorf("token is invalid (exp is not string)")
	// }
	// expTime := time.Unix(expTimeUNIX, 0)

	// if time.Now().After(expTime) {
	// 	return fmt.Errorf("token is expired")
	//}

	return claims, nil
}

func Field(field, tokenStringBearer, secret string) (interface{}, error) {
	claims, err := Verify(tokenStringBearer, secret)
	if err != nil {
		return nil, err
	}

	fieldVal, ok := claims[field]
	if !ok {
		return nil, errors.New("field is not present")
	}

	return fieldVal, nil
}
