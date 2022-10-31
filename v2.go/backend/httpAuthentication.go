// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

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
	authorized, userUuid, firstName, lastName := isDbAuthorized(dbName, login.Id, login.Password)
	if !authorized || userUuid == "" {
		c.JSON(http.StatusBadRequest, "")
		return
	}

	tokenString, err := getNewToken(dbName, login.Id, userUuid, firstName, lastName)
	if err != nil {
		utils.LogError("Failed to generate signed string.")
		c.JSON(http.StatusServiceUnavailable, "")
	}

	utils.Log("Token generated: %v.", tokenString)
	jwtToken := JWTtoken{tokenString}
	c.JSON(http.StatusOK, jwtToken)
}

func getNewToken(dbName string, id string, userUuid string, firstName string, lastName string) (string, error) {
	expiration := time.Now().Add(time.Duration(utils.Configuration.HttpServer.JwtExpiration) * time.Minute)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":          id,
		"userUuid":      userUuid,
		"userFirstName": firstName,
		"userLastName":  lastName,
		"expires":       expiration,
	})
	utils.Log("Token generated for %v, expiration: %v", id, expiration)
	jwtSecret := utils.GetJWTSecret(dbName)
	return token.SignedString([]byte(jwtSecret))
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
			c.Set("authorized", false)
			utils.LogError("No authorization found in header for request %v", c.Request)
			return
		}

		if len(header) < 10 {
			c.Set("authorized", false)
			utils.LogError("Incorrect header for request %v", c.Request)
			return
		}
		var tokenString = header[7:]

		jwtSecret := utils.GetJWTSecret(dbName)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})
		if token == nil {
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			user := claims["user"]
			today := time.Now()
			expiration := claims["expires"]
			expirationDate, _ := time.Parse(time.RFC3339Nano, fmt.Sprintf("%v", expiration))

			if today.After(expirationDate) {
				c.Set("authorized", false)
				utils.Log("[%v] Authorization expired (%v).", user, expirationDate)
				c.Abort()
				c.IndentedJSON(http.StatusUnauthorized,
					gin.H{
						"error":      "Authorization expired.",
						"disconnect": true})
				return
			}
			c.Set("authorized", true)
		} else {
			utils.LogError("Invalid request: %v.", err)
			c.Set("authorized", false)
			c.Abort()
			c.IndentedJSON(http.StatusUnauthorized,
				gin.H{
					"error":      "Unauthorized.",
					"disconnect": true})
			return
		}
	}
}
