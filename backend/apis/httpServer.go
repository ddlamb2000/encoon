// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2025

package apis

import (
	"net/http"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/model"
	"github.com/gin-gonic/gin"
)

func SetApiRoutes(r *gin.Engine) {
	v1 := r.Group("/:dbName/api/v1")
	{
		v1.POST("/authentication", authentication)
		v1.GET("/", authMiddleware(), GetGridsRowsApi)
		v1.GET("/:gridUuid", authMiddleware(), GetGridsRowsApi)
		v1.GET("/:gridUuid/:uuid", authMiddleware(), GetGridsRowsApi)
		v1.POST("/:gridUuid", authMiddleware(), PostGridsRowsApi)
		v1.POST("/:gridUuid/:uuid", authMiddleware(), PostGridsRowsApi)
	}
}

func getParameters(c *gin.Context) (ApiParameters, error) {
	dbName, gridUuid, uuid := c.Param("dbName"), c.Param("gridUuid"), c.Param("uuid")
	userUuid, userName := c.GetString("userUuid"), c.GetString("user")
	auth, exists := c.Get("authorized")
	p := ApiParameters{
		DbName:   dbName,
		UserUuid: userUuid,
		UserName: userName,
		GridUuid: gridUuid,
		Uuid:     uuid,
	}
	if dbName == "" || gridUuid == "" || len(userUuid) != len(model.UuidUsers) || userName == "" || !exists || auth == false {
		return p, configuration.LogAndReturnError(dbName, userName, "User not authorized.")
	}
	p.filterColumnOwned = c.Query("filterColumnOwned") == "true"
	p.filterColumnName = c.Query("filterColumnName")
	p.filterColumnGridUuid = c.Query("filterColumnGridUuid")
	p.filterColumnValue = c.Query("filterColumnValue")
	return p, nil
}

func getHttpErrorCode(response GridResponse) int {
	switch {
	case response.System:
		return http.StatusInternalServerError
	case response.Forbidden:
		return http.StatusForbidden
	case response.TimeOut:
		return http.StatusRequestTimeout
	default:
		return http.StatusNotFound
	}
}

func GetGridsRowsApi(c *gin.Context) {
	p, err := getParameters(c)
	if err != nil {
		c.Abort()
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	response := GetGridsRows(c.Request.Context(), c.Request.RequestURI, p, GridPost{})
	if response.Err != nil {
		c.Abort()
		c.JSON(getHttpErrorCode(response), gin.H{"error": response.Err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"response": response})
}

func PostGridsRowsApi(c *gin.Context) {
	p, err := getParameters(c)
	if err != nil {
		c.Abort()
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}
	var payload GridPost
	c.ShouldBindJSON(&payload)
	response := PostGridsRows(c.Request.Context(), c.Request.RequestURI, p, payload)
	if response.Err != nil {
		c.Abort()
		c.JSON(getHttpErrorCode(response), gin.H{"error": response.Err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"response": response})
}
