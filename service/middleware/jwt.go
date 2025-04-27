package jwt

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/singhpranshu/splitbill/common"
	"github.com/singhpranshu/splitbill/repository/model"
)

type UserClaims struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

var secretKey []byte

type UserClaimKey string

const UserClaimKeyName UserClaimKey = "userClaim"

func InitJwt(c *common.Config) {
	// Initialize the JWT middleware with the secret key
	secretKey = []byte(c.JWTSecretKey)
}

func CreateToken(user model.User) (string, error) {
	// Create a new JWT token with the payload
	userClaim := UserClaims{
		UserID: user.ID,
		Name:   user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "splitbill",
			Subject:   user.Username,
			Audience:  []string{"splitbill"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        fmt.Sprintf("%d", user.ID),
		},
		// Add any other custom claims you need
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim)
	// Sign the token using the secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		// Token is valid, return the claims
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func AuthenticateMiddleware(originalHandler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authTokenCookie, err := r.Cookie("token")
		if err != nil || authTokenCookie.Value == "" {
			log.Println("Authorization header is required", err, authTokenCookie)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(common.GetHttpErrorResponse(http.StatusUnauthorized, "Authorization header is required")))
			return
		}
		userClaim, err := VerifyToken(authTokenCookie.Value)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(common.GetHttpErrorResponse(http.StatusUnauthorized, "Invalid token")))
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserClaimKeyName, userClaim)
		originalHandler.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthenticateMiddlewareHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		AuthenticateMiddleware(next.ServeHTTP)(w, r)
	})
}
