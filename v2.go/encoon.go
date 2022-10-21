// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"d.lambert.fr/encoon/core"
	"d.lambert.fr/encoon/utils"
	"github.com/gin-gonic/gin"
)

const httpPort = ":8080"

func setAndStartServer() *http.Server {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*.html")
	router.Static("/stylesheets", "./stylesheets")
	router.Static("/javascript", "./javascript")
	router.Static("/images", "./images")
	router.StaticFile("favicon.ico", "./images/favicon.ico")

	router.GET("/", core.GetIndexHtml)
	router.GET("/users.html", core.GetUsersHtml)

	router.GET("/ping", core.PingApi)
	v1 := router.Group("/api/v1")
	{
		v1.GET("/users", core.GetUsersApi)
		v1.GET("/users/:uuid", core.GetUserByIDApi)
		v1.POST("/users", core.PostUsersApi)
	}

	srv := &http.Server{
		Addr:         httpPort,
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

func initServers() *http.Server {
	initWithLog()
	srv := setAndStartServer()
	core.LoadData()
	return srv
}

const (
	dbHost     = "localhost"
	dbPort     = 5432
	dbUser     = "david.lambert"
	dbPassword = ""
	dbName     = "david.lambert"
)

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	defer db.Close()

	srv := initServers()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	shutDownServer(srv, ctx)
}
