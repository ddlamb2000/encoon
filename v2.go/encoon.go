// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package main

import (
	"io"
	"os"

	"d.lambert.fr/encoon/core"
	"github.com/gin-gonic/gin"
)

func setRouter() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*") // see templates https://pkg.go.dev/text/template, https://gohugo.io/templates/
	router.Static("/stylesheets", "./stylesheets")
	router.GET("/users", core.GetUsersJson)
	router.GET("/users/:uuid", core.GetAlbumByIDJson)
	router.POST("/users", core.PostUsers)
	router.GET("/users.html", core.GetUsersHtml)
	router.GET("/ping", core.Ping)
	return router
}

func main() {
	f, _ := os.Create("logs/encoon.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	core.LoadData()
	router := setRouter()
	router.Run("localhost:8080")
}
