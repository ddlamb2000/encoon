// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package middleware

import (
	"fmt"
	"net/http"
	"time"

	"d.lambert.fr/encoon/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type login struct {
	Id       string `json:"id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type JWTtoken struct {
	Token string `json:"token" binding:"required"`
}

func authentication(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	dbName := c.Param("dbName")
	if dbName == "" {
		c.JSON(http.StatusBadRequest, "")
		return
	}

	var login login
	c.BindJSON(&login)
	authorized, userUuid := isDbAuthorized(dbName, login.Id, login.Password)
	if !authorized || userUuid == "" {
		c.JSON(http.StatusBadRequest, "")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":       login.Id,
		"userUuid":   userUuid,
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
		if len(header) < 10 {
			utils.LogError("Incorrect header for request %v", c.Request)
			return
		}
		var tokenString = header[7:]
		utils.Log("Got token string.")

		jwtSecret := utils.GetJWTSecret(dbName)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		utils.Log("Extracting claims.")
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			user := claims["user"]
			userUuid := claims["userUuid"]
			today := time.Now()
			expiration := claims["expiration"]
			expirationDate, _ := time.Parse(time.RFC3339Nano, fmt.Sprintf("%v", expiration))
			utils.Log("User: %v, userUuid: %v, expiration: %v.", user, userUuid, expirationDate)
			if today.After(expirationDate) {
				utils.Log("Token expired.")
				c.Abort()
				c.IndentedJSON(http.StatusUnauthorized,
					gin.H{
						"error":      true,
						"message":    "Token expired.",
						"disconnect": true})
				return
			}
		} else {
			utils.LogError("Invalid request: %v.", err)
			c.Abort()
			c.IndentedJSON(http.StatusUnauthorized,
				gin.H{
					"error":      true,
					"message":    "Unauthorized.",
					"disconnect": true})
			return
		}
	}
}
