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
)

var (
	router = gin.Default()
	srv    *http.Server
)

func SetAndStartHttpServer() {
	setHtmlTemplates()
	setStaticFiles()
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

func setHtmlTemplates() {
	router.LoadHTMLGlob("frontend/templates/*.html")
}

func setStaticFiles() {
	router.Static("/stylesheets", "./frontend/stylesheets")
	router.Static("/javascript", "./frontend/javascript")
	router.Static("/images", "./frontend/images")
	router.StaticFile("favicon.ico", "./frontend/images/favicon.ico")
}

func setHtmlRoutes() {
	router.GET("/", getIndexHtml)
	router.GET("/:dbName", getIndexHtml)
	router.GET("/:dbName/users", getIndexHtml)
	router.GET("/:dbName/users/:uuid", getIndexHtml)
}

func setApiRoutes() {
	v1 := router.Group("/:dbName/api/v1")
	{
		v1.POST("/authentication", authentication)
		v1.GET("/users", authMiddleware(), core.GetUsersApi)
		v1.GET("/users/:uuid", authMiddleware(), core.GetUserByIDApi)
		v1.POST("/users", authMiddleware(), core.PostUsersApi)
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
	uuid := c.Param("uuid")
	if dbName == "" {
		c.HTML(http.StatusOK, "home.html", gin.H{"appName": "εncooη"})
	} else if utils.IsDatabaseEnabled(dbName) {
		c.HTML(http.StatusOK, "index.html", gin.H{"appName": "εncooη", "dbName": dbName, "uuid": uuid})
	} else {
		c.HTML(http.StatusNotFound, "nofound.html", gin.H{"appName": "εncooη"})
	}
}
