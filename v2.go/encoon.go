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

func setRouter() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html") // see templates https://pkg.go.dev/text/template, https://gohugo.io/templates/
	router.Static("/stylesheets", "./stylesheets")
	router.GET("/users.html", core.GetUsersHtml)
	router.GET("/ping", core.Ping)
	v1 := router.Group("/v1")
	{
		v1.GET("/users", core.GetUsersJson)
		v1.GET("/users/:uuid", core.GetAlbumByIDJson)
		v1.POST("/users", core.PostUsers)
	}
	return router
}

func main() {
	f, _ := os.Create("logs/encoon.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	utils.Log("Starting.")
	core.LoadData()
	router := setRouter()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.LogFatal("Listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 2 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	utils.Log("Shuting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		utils.LogFatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 2 seconds.
	select {
	case <-ctx.Done():
		utils.Log("Timeout of 2 seconds.")
	}
	utils.Log("Server exiting")
}
