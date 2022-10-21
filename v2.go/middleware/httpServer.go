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
		utils.LogFatal("Listen:", err)
		return
	}
}

func setHtmlRoutes() {
	router.LoadHTMLGlob("frontend/templates/*.html")
	router.Static("/stylesheets", "./frontend/stylesheets")
	router.Static("/javascript", "./frontend/javascript")
	router.Static("/images", "./frontend/images")
	router.StaticFile("favicon.ico", "./frontend/images/favicon.ico")
	router.GET("/", getHomeHtml)
	router.GET("/:db", getIndexHtml)
}

func setApiRoutes() {
	v1 := router.Group("/:db/api/v1")
	{
		v1.GET("/users", core.GetUsersApi)
		v1.GET("/users/:uuid", core.GetUserByIDApi)
		v1.POST("/users", core.PostUsersApi)
	}
}

func ShutDownHttpServer(ctx context.Context) {
	if err := srv.Shutdown(ctx); err != nil {
		utils.LogFatal("Server Shutdown:", err)
		return
	}
	utils.Log("Http server stopped.")
}

func getHomeHtml(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", gin.H{"title": "εncooη --- no database "})
}

func getIndexHtml(c *gin.Context) {
	db := c.Param("db")
	if utils.DatabaseAllowed(db) {
		c.HTML(http.StatusOK, "index.html", gin.H{"title": "εncooη", "db": db})
	} else {
		c.HTML(http.StatusForbidden, "forbidden.html", gin.H{"title": "εncooη"})
	}
}
