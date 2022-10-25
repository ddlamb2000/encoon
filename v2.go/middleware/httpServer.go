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
	router.GET("/:db/users", getUsersHtml)
	router.GET("/:db/user/:uuid", getUserHtml)
	router.GET("/:db/login", getLoginHtml)
	router.GET("/:db", getIndexHtml)
}

func setApiRoutes() {
	v1 := router.Group("/:db/api/v1")
	{
		v1.GET("/users", core.GetUsersApi)
		v1.GET("/user/:uuid", core.GetUserByIDApi)
		v1.POST("/users", core.PostUsersApi)
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
	db := c.Param("db")
	if db == "" {
		c.HTML(http.StatusOK, "home.html", gin.H{"title": "εncooη --- no database "})
	} else if utils.IsDatabaseEnabled(db) {
		c.HTML(http.StatusOK, "index.html", gin.H{"title": "εncooη", "db": db})
	} else {
		c.HTML(http.StatusNotFound, "nofound.html", gin.H{"title": "εncooη"})
	}
}

func getLoginHtml(c *gin.Context) {
	db := c.Param("db")
	if utils.IsDatabaseEnabled(db) {
		c.HTML(http.StatusOK, "login.html", gin.H{"title": "εncooη", "db": db})
	} else {
		c.HTML(http.StatusNotFound, "nofound.html", gin.H{"title": "εncooη"})
	}
}

func getUsersHtml(c *gin.Context) {
	db := c.Param("db")
	if utils.IsDatabaseEnabled(db) {
		token, _ := generateJWT(db)
		c.Header("Token", token)
		c.HTML(http.StatusOK, "users.html", gin.H{"title": "εncooη", "db": db})
	} else {
		c.HTML(http.StatusNotFound, "nofound.html", gin.H{"title": "εncooη"})
	}
}

func getUserHtml(c *gin.Context) {
	db := c.Param("db")
	uuid := c.Param("uuid")
	if utils.IsDatabaseEnabled(db) {
		c.HTML(http.StatusOK, "users.html", gin.H{"title": "εncooη", "db": db, "uuid": uuid})
	} else {
		c.HTML(http.StatusNotFound, "nofound.html", gin.H{"title": "εncooη"})
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
