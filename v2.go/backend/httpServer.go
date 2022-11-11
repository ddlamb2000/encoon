// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import (
	"d.lambert.fr/encoon/apis"
	"d.lambert.fr/encoon/authentication"
	"github.com/gin-gonic/gin"
)

func SetApiRoutes(r *gin.Engine) {
	v1 := r.Group("/:dbName/api/v1")
	{
		v1.POST("/authentication", authentication.Authentication)
		v1.GET("/", authentication.AuthMiddleware(), apis.GetGridsRowsApi)
		v1.GET("/:gridUri", authentication.AuthMiddleware(), apis.GetGridsRowsApi)
		v1.GET("/:gridUri/:uuid", authentication.AuthMiddleware(), apis.GetGridsRowsApi)
		v1.POST("/:gridUri", authentication.AuthMiddleware(), apis.PostGridsRowsApi)
	}
}
