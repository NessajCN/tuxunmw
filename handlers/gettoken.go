package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

const JWTSECRET string = "m4+gHOSXrWtTkzzblR8fVupaJkeMaXKlrtpAGmwjTWw="

// Verify jwt token string. Return the payload username if verified.
func verifyJWT(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &usernameClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(JWTSECRET), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*usernameClaims); ok && token.Valid {
		return claims.Name, nil
	} else {
		return "", errors.New("invalid token")
	}
}

// userAuth returns username
// if authorized.
func userAuth(tokenString string, conf config) (string, error) {
	tokenuser, parseerr := verifyJWT(tokenString)
	if parseerr != nil {
		return "", parseerr
	}
	if tokenuser != conf.Username {
		return "", errors.New("username unauthenticated")
	}
	return tokenuser, nil
}

func GetToken(ctx *gin.Context) {
	var tokenauth tokenAuthReq
	if err := ctx.BindJSON(&tokenauth); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"success": false,
			"payload": tokenauth,
		})
		return
	}
	var conf config
	if err := parseConfig(&conf); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"success": false,
		})
		return
	}

	if tokenauth.Name != conf.Username || tokenauth.Password != conf.Password {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
			"success": false,
		})
		return
	}

	claims := usernameClaims{
		tokenauth.Name,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 7 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenstring, err := token.SignedString([]byte(JWTSECRET))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"success": false,
		})
		return
	}
	tokenreply := accesstoken{
		Message: "Token generated",
		Success: true,
		Token:   tokenstring,
	}
	ctx.JSON(http.StatusOK, tokenreply)
}
