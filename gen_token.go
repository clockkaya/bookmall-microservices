package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func main() {
	// 必须和你 search-api.yaml 里的 AccessSecret 一模一样！
	secret := "u-really-need-to-change-this-secret"

	now := time.Now().Unix()
	claims := make(jwt.MapClaims)
	claims["exp"] = now + 86400
	claims["iat"] = now
	// 模拟 userId = 8888
	claims["userId"] = 1

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	tokenString, _ := token.SignedString([]byte(secret))

	fmt.Println("Bearer " + tokenString)
}
