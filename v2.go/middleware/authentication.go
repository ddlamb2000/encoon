// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package middleware

import (
	"fmt"
	"net/http"
	"time"

	"d.lambert.fr/encoon/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type LOGIN struct {
	ID       string `json:"id" binding:"required"`
	PASSWORD string `json:"password" binding:"required"`
}

type JWTtoken struct {
	TOKEN string `json:"token" binding:"required"`
}

var authID = "david.lambert"
var authPassword = "hello?"

func authentication(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	dbName := c.Param("dbName")
	if dbName == "" {
		c.JSON(http.StatusBadRequest, "")
		return
	}

	var login LOGIN
	c.BindJSON(&login)

	if authID != login.ID {
		c.JSON(http.StatusBadRequest, "")
		return
	}
	if authPassword != login.PASSWORD {
		c.JSON(http.StatusBadRequest, "")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":       login.ID,
		"timestamp":  int32(time.Now().Unix()),
		"expiration": time.Now().Add(10 * time.Minute),
	})
	jwtSecret := utils.GetJWTSecret(dbName)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		utils.LogError("Failed to generate signed string.")
		c.JSON(http.StatusServiceUnavailable, "")
	}

	utils.Log("Token generated: %v.", tokenString)
	jwtToken := JWTtoken{tokenString}
	c.JSON(http.StatusOK, jwtToken)
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbName := c.Param("dbName")
		if dbName == "" {
			utils.LogError("No database name for request %v", c.Request)
			return
		}

		var header = c.Request.Header.Get("Authorization")
		if header == "" {
			utils.LogError("No authorization found in header for request %v", c.Request)
			return
		}

		utils.Log("Got authorization.")
		var tokenString = header[7:]
		utils.Log("Got token string: %v.", tokenString)

		jwtSecret := utils.GetJWTSecret(dbName)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		utils.Log("Extracting claims.")
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			utils.Log("User: %v, timestamp: %v, expiration: %v.", claims["user"], claims["timestamp"], claims["expiration"])
			today := time.Now()
			expiration := claims["expiration"]
			expirationDate, _ := time.Parse(time.RFC3339Nano, fmt.Sprintf("%v", expiration))
			if today.After(expirationDate) {
				c.Abort()
				c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": true, "message": "Expired."})
				return
			}

		} else {
			utils.LogError("Invalid request: %v.", err)
			c.Abort()
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Writer.Write([]byte("Unauthorized"))
			return
		}
	}
}
