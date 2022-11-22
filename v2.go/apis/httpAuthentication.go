// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"fmt"
	"net/http"
	"time"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/database"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type login struct {
	Id       string `json:"id" yaml:"id" binding:"required"`
	Password string `json:"password" yaml:"password" binding:"required"`
}

type JWTtoken struct {
	Token string `json:"token" yaml:"token" binding:"required"`
}

func authentication(c *gin.Context) {
	dbName := c.Param("dbName")
	if dbName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No database parameter"})
		return
	}
	var login login
	c.ShouldBindJSON(&login)
	userUuid, firstName, lastName, timeOut, err := database.IsDbAuthorized(c.Request.Context(), dbName, login.Id, login.Password)
	if err != nil || userUuid == "" {
		c.Abort()
		if timeOut {
			c.JSON(http.StatusRequestTimeout, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}
		return
	}
	expiration := time.Now().Add(time.Duration(configuration.GetConfiguration().HttpServer.JwtExpiration) * time.Minute)
	tokenString, err := getNewToken(dbName, login.Id, userUuid, firstName, lastName, expiration)
	if err != nil {
		c.Abort()
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}
	logUri(c, dbName, login.Id)
	c.JSON(http.StatusOK, JWTtoken{tokenString})
}

func getNewToken(dbName string, user string, userUuid string, firstName string, lastName string, expiration time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":          user,
		"userUuid":      userUuid,
		"userFirstName": firstName,
		"userLastName":  lastName,
		"expires":       expiration,
	})
	configuration.Log(dbName, user, "Token generated, expiration: %v", expiration)
	jwtSecret := configuration.GetJWTSecret(dbName)
	return token.SignedString([]byte(jwtSecret))
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbName := c.Param("dbName")
		if dbName == "" {
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No database parameter."})
			return
		}

		var header = c.Request.Header.Get("Authorization")
		if header == "" {
			c.Set("authorized", false)
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization found."})
			return
		}

		if len(header) < 10 {
			c.Set("authorized", false)
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect header."})
			return
		}
		var tokenString = header[7:]

		jwtSecret := configuration.GetJWTSecret(dbName)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, configuration.LogAndReturnError(dbName, "", "Unexpect signing method: %v.", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})
		if token == nil {
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			user := claims["user"]
			userUuid := claims["userUuid"]
			today := time.Now()
			expiration := claims["expires"]
			expirationDate, _ := time.Parse(time.RFC3339Nano, fmt.Sprintf("%v", expiration))
			userName := fmt.Sprintf("%v", user)

			if today.After(expirationDate) {
				c.Set("authorized", false)
				configuration.Log(dbName, userName, "Authorization expired (%v).", expirationDate)
				c.Abort()
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization expired.", "expired": true})
				return
			}
			c.Set("authorized", true)
			c.Set("user", user)
			c.Set("userUuid", userUuid)
			logUri(c, dbName, user)
		} else {
			configuration.LogError(dbName, "", "Invalid request: %v.", err)
			c.Set("authorized", false)
			c.Abort()
			c.JSON(http.StatusUnauthorized,
				gin.H{"error": fmt.Sprintf("Invalid request or unauthorized database access: %v.", err)})
			return
		}
	}
}

func logUri(c *gin.Context, dbName string, user any) {
	userName := fmt.Sprintf("%v", user)
	configuration.Log(dbName, userName, "%v", c.Request.RequestURI)
}
