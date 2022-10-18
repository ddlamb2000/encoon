// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package main

import (
	"fmt"
	"io"
	"os"

	"d.lambert.fr/encoon/core"
	"github.com/gin-gonic/gin"
)

func main() {
	f, _ := os.Create("logs/encoon.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	fmt.Fprintf(gin.DefaultWriter, "[encoon] εncooη : data structuration, presentation and navigation.\n")

	core.LoadUsers()

	router := gin.Default()
	router.LoadHTMLGlob("templates/*") // see templates https://pkg.go.dev/text/template, https://gohugo.io/templates/
	router.Static("/stylesheets", "./stylesheets")
	router.GET("/albums.html", core.GetAlbumsHtml)
	router.GET("/albums", core.GetAlbums)
	router.GET("/albums/:id", core.GetAlbumByID)
	router.POST("/albums", core.PostAlbums)

	router.GET("/users", core.GetUsersJson)

	router.Run("localhost:8080")
}
