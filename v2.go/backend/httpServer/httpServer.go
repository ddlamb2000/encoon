// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package httpServer

import (
	"context"
	"net/http"
	"time"

	"d.lambert.fr/encoon/backend/core"
	"d.lambert.fr/encoon/backend/utils"
	"github.com/gin-gonic/gin"
)

const httpPort = ":8080"

var (
	srv    *http.Server
	router *gin.Engine
)

func SetAndStartServer() {
	router = gin.Default()
	setHtmlRoutes()
	setApiRoutes()
	srv = &http.Server{
		Addr:         httpPort,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.LogFatal("Listen:", err)
			return
		}
	}()
	utils.Log("Listen on port " + httpPort)
}

func setHtmlRoutes() {
	router.LoadHTMLGlob("frontend/templates/*.html")
	router.Static("/stylesheets", "./frontend/stylesheets")
	router.Static("/javascript", "./frontend/javascript")
	router.Static("/images", "./frontend/images")
	router.StaticFile("favicon.ico", "./images/favicon.ico")

	router.GET("/", core.GetIndexHtml)
	router.GET("/users.html", core.GetUsersHtml)
}

func setApiRoutes() {
	router.GET("/ping", core.PingApi)
	v1 := router.Group("/api/v1")
	{
		v1.GET("/users", core.GetUsersApi)
		v1.GET("/users/:uuid", core.GetUserByIDApi)
		v1.POST("/users", core.PostUsersApi)
	}
}

func ShutDownServer(ctx context.Context) {
	utils.Log("Shut down server on port " + httpPort + ".")
	if err := srv.Shutdown(ctx); err != nil {
		utils.LogFatal("Server Shutdown:", err)
		return
	}
	utils.Log("Server on port " + httpPort + " stopped.")
}
