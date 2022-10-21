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

func ShutDownHttpServer(ctx context.Context) {
	utils.Log("Stopping http server.")
	if err := srv.Shutdown(ctx); err != nil {
		utils.LogFatal("Server Shutdown:", err)
		return
	}
	utils.Log("Http server stopped.")
}
