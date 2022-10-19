// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"d.lambert.fr/encoon/core"
	"d.lambert.fr/encoon/utils"
	"github.com/gin-gonic/gin"
)

const portHtml = ":8080"
const portApi = ":8081"

func setAndStartServerHtml() *http.Server {
	router := gin.Default()
	// see https://pkg.go.dev/text/template,
	// see https://gohugo.io/templates/
	router.LoadHTMLGlob("templates/*.html")
	router.Static("/stylesheets", "./stylesheets")
	router.GET("/users.html", core.GetUsersHtml)
	srv := &http.Server{
		Addr:         portHtml,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	startServer(srv)
	return srv
}

func setAndStartServerApi() *http.Server {
	router := gin.Default()
	router.GET("/ping", core.PingApi)
	v1 := router.Group("/v1")
	{
		v1.GET("/users", core.GetUsersApi)
		v1.GET("/users/:uuid", core.GetUserByIDApi)
		v1.POST("/users", core.PostUsersApi)
	}
	srv := &http.Server{
		Addr:         portApi,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	startServer(srv)
	return srv
}

func initWithLog() {
	f, _ := os.Create("logs/encoon.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	utils.Log("Starting.")
}

func startServer(srv *http.Server) {
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.LogFatal("Listen:", err)
			return
		}
	}()
	utils.Log("Listen on port " + srv.Addr)
}

func shutDownServer(srv *http.Server, ctx context.Context) {
	utils.Log("Shut down server on port " + srv.Addr + ".")
	if err := srv.Shutdown(ctx); err != nil {
		utils.LogFatal("Server Shutdown:", err)
		return
	}
	utils.Log("Server on port " + srv.Addr + " stopped.")
}

func initServers() (*http.Server, *http.Server) {
	initWithLog()
	srvHtml := setAndStartServerHtml()
	srvApi := setAndStartServerApi()
	core.LoadData()
	return srvHtml, srvApi
}

func main() {
	srvHtml, srvApi := initServers()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	shutDownServer(srvHtml, ctx)
	shutDownServer(srvApi, ctx)
}
