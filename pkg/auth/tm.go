package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"log"
	"time"
)

var singingKey = viper.GetString("key.jwt")

type TokenManager struct {
	db *pgxpool.Pool
}

func Error(err error) {
	if err != nil {
		log.Println(err.Error())
		return
	}
}

type TokenClaims struct {
	*jwt.RegisteredClaims
	UserId int `json:"user_id"`
}

func NewTokenManager(conn *pgxpool.Pool) TokenManager {
	return TokenManager{db: conn}
}

func (t *TokenManager) GenerateToken(id int) (string, error) {
	if id == 0 {
		return "", errors.New("error: id = 0; unauthorized")
	}

	issuedAt := jwt.NewNumericDate(time.Now())
	expiresAccess := jwt.NewNumericDate(time.Now().Add(60 * 24 * 365 * time.Minute))

	accessClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		&jwt.RegisteredClaims{
			IssuedAt:  issuedAt,
			ExpiresAt: expiresAccess,
		},
		id,
	})
	//Gen Access Token
	accessToken, err := accessClaims.SignedString([]byte(singingKey))
	Error(err)

	return accessToken, nil
}

func (t *TokenManager) ValidateToken(accessToken string) (int, error) {
	//var valid bool

	aToken, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid access-token")
		}
		return []byte(singingKey), nil
	})
	if err != nil {
		log.Print(err.Error())
	}

	claims, ok := aToken.Claims.(*TokenClaims)
	if !ok {
		return 0, errors.New("invalid token") //TODO
	}
	return claims.UserId, nil
}
