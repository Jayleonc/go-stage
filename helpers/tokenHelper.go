package helpers

import (
	"fmt"
	"github.com/Jayleonc/go-stage/database"
	"github.com/Jayleonc/go-stage/models"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"os"
	"time"
)

type SignedDetails struct {
	Email    string
	UserName string
	UserId   string
	UserType string
	jwt.RegisteredClaims
}

var SecretKey = os.Getenv("SECRET_KEY")

var db = database.DBInstance()

func GenerateAllTokens(email string, username string, userType string, userId string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:    email,
		UserName: username,
		UserId:   userId,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * time.Duration(24))),
		},
	}

	refreshClaims := &SignedDetails{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * time.Duration(24))),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SecretKey))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SecretKey))
	if err != nil {
		fmt.Println("GenerateAllTokens() error")
		log.Panic(err)
		return
	}
	return token, refreshToken, err
}
func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
	var user models.User
	user.Token = signedToken
	user.RefreshToken = signedRefreshToken
	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UpdatedAt = updatedAt
	save := db.Debug().Model(&user).Where("user_id = ?", userId).Updates(user)
	if save.RowsAffected == 0 {
		var msg = fmt.Sprintf("更新 token 失败")
		log.Panic(msg)
		return
	}
	return
}
func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		},
	)
	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("the token is invalid")
		msg = err.Error()
		return
	}
	nowDate := jwt.NewNumericDate(time.Now().Local())
	if claims.ExpiresAt.Unix() < nowDate.Unix() {
		msg = fmt.Sprintf("token is expired")
		msg = err.Error()
		return
	}
	return claims, msg
}
