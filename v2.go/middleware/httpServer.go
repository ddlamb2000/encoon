// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"d.lambert.fr/encoon/backend/core"
	"d.lambert.fr/encoon/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var (
	router = gin.Default()
	srv    *http.Server
)

func SetAndStartHttpServer() {
	setHtmlRoutes()
	setApiRoutes()
	srv = &http.Server{
		Addr:         fmt.Sprintf(":%d", utils.Configuration.HttpServer.Port),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	utils.Log("Listening http.")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		utils.LogError("Error on http listening: %v.", err)
		return
	}
}

func setHtmlRoutes() {
	router.LoadHTMLGlob("frontend/templates/*.html")
	router.Static("/stylesheets", "./frontend/stylesheets")
	router.Static("/javascript", "./frontend/javascript")
	router.Static("/images", "./frontend/images")
	router.StaticFile("favicon.ico", "./frontend/images/favicon.ico")
	router.GET("/", getIndexHtml)
	router.GET("/:dbName", getIndexHtml)
	router.GET("/:dbName/users", getUsersHtml)
	router.GET("/:dbName/user/:uuid", getUserHtml)
}

func setApiRoutes() {
	v1 := router.Group("/:dbName/api/v1")
	{
		v1.GET("/users", core.GetUsersApi)
		v1.POST("/users", core.PostUsersApi)
		v1.GET("/user/:uuid", core.GetUserByIDApi)
	}
}

func ShutDownHttpServer(ctx context.Context) {
	if err := srv.Shutdown(ctx); err != nil {
		utils.LogError("Server shutdown: %v.", err)
		return
	}
	utils.Log("Http server stopped.")
}

func getIndexHtml(c *gin.Context) {
	dbName := c.Param("dbName")
	if dbName == "" {
		c.HTML(http.StatusOK, "home.html", gin.H{"title": "εncooη"})
	} else if utils.IsDatabaseEnabled(dbName) {
		c.HTML(http.StatusOK, "index.html", gin.H{"appName": "εncooη", "dbName": dbName})
	} else {
		c.HTML(http.StatusNotFound, "nofound.html", gin.H{"appName": "εncooη"})
	}
}

func getUsersHtml(c *gin.Context) {
	dbName := c.Param("dbName")
	if utils.IsDatabaseEnabled(dbName) {
		token, _ := generateJWT(dbName)
		c.Header("Token", token)
		c.HTML(http.StatusOK, "users.html", gin.H{"appName": "εncooη", "dbName": dbName})
	} else {
		c.HTML(http.StatusNotFound, "nofound.html", gin.H{"appName": "εncooη"})
	}
}

func getUserHtml(c *gin.Context) {
	dbName := c.Param("dbName")
	uuid := c.Param("uuid")
	if utils.IsDatabaseEnabled(dbName) {
		c.HTML(http.StatusOK, "users.html", gin.H{"appName": "εncooη", "dbName": dbName, "uuid": uuid})
	} else {
		c.HTML(http.StatusNotFound, "nofound.html", gin.H{"appName": "εncooη"})
	}
}

func generateJWT(dbName string) (string, error) {
	jwtSecret := utils.GetJWTSecret(dbName)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(10 * time.Minute)
	claims["authorized"] = true
	claims["user"] = "username"
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		utils.LogError("Signing Error: %v.", err)
		return "Signing Error", err
	}
	return tokenString, nil
}
