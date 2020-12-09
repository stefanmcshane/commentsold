package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type Auth struct {
	Username string
	Password string
}

type SignedToken struct {
	Username string
	Token    string
}

type JWTToken struct {
	SigningKey string
}

// GenerateToken returns a JWT token for a given Auth user, with a 15 minute timeout
func (jt JWTToken) GenerateToken(ctx context.Context, au Auth) (*SignedToken, error) {

	if au.Username == "" {
		e := "No username provided"
		lg(ctx).Error(e)
		return nil, errors.New(e)
	}

	lg(ctx).Infof("Generating token for user %s", au.Username)
	timeout := 15
	expirationTime := time.Now().Add(time.Duration(timeout) * time.Minute)

	claims := jwt.MapClaims{
		"username": au.Username,
		"claims": jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jt.SigningKey))
	if err != nil {
		lg(ctx).Errorf("Error signing token - %s", err.Error())
		return nil, err
	}

	lg(ctx).Infof("Token successfully generated for %s", au.Username)
	t := &SignedToken{
		Username: au.Username,
		Token:    signedToken,
	}

	return t, nil
}

// ValidateToken takes a SignedToken, and ensures its valid against the signing key.
// If valid, the username attached to the token will be returned
func (jt JWTToken) ValidateToken(ctx context.Context, st SignedToken) (string, error) {
	signedToken := st.Token
	if signedToken == "" {
		e := "Token was empty"
		lg(ctx).Error(e)
		return "", errors.New(e)
	}

	signedToken = strings.ReplaceAll(signedToken, "Bearer ", "")

	// Check token for structural validity
	cl := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(signedToken, cl, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			e := fmt.Sprintf("Unexpected signing method: %v. Potential Security Issue - This token has been changed externally", token.Header["alg"])
			lg(ctx).Error(e)
			return nil, errors.New(e)
		}
		return []byte(jt.SigningKey), nil
	})
	if err != nil {
		e := fmt.Sprintf("JWT Token parsing error - %s", err.Error())
		lg(ctx).Error(e)
		return "", errors.New(e)
	}
	if !token.Valid {
		e := "Invalid token. Potential security issue"
		lg(ctx).Error(e)
		return "", errors.New(e)
	}

	claims, ok := cl["claims"].(map[string]interface{})
	if !ok {
		e := "Unable to retrieve claims from token"
		lg(ctx).Error(e)
		return "", errors.New(e)
	}

	// Check token for expected contents - username and time
	expiry, ok := claims["exp"].(float64)
	if !ok {
		e := "Expiry time missing from token or malformed"
		lg(ctx).Error(e)
		return "", errors.New(e)
	}

	if time.Now().Sub(time.Unix(int64(expiry), 0)).Minutes() > 15 {
		e := "Token has expired"
		lg(ctx).Error(e)
		return "", errors.New(e)
	}

	username, ok := cl["username"].(string)
	if !ok {
		e := "Username missing from token"
		lg(ctx).Error(e)
		return "", errors.New(e)
	}

	return username, nil

}
