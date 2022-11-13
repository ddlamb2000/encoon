// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"github.com/gin-gonic/gin"
)

func SetApiRoutes(r *gin.Engine) {
	v1 := r.Group("/:dbName/api/v1")
	{
		v1.POST("/authentication", authentication)
		v1.GET("/", authMiddleware(), GetGridsRowsApi)
		v1.GET("/:gridUri", authMiddleware(), GetGridsRowsApi)
		v1.GET("/:gridUri/:uuid", authMiddleware(), GetGridsRowsApi)
		v1.POST("/:gridUri", authMiddleware(), PostGridsRowsApi)
	}
}
