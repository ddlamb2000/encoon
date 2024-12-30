// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package apis

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/database"
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

func (r ApiRequest) log(format string, a ...any) {
	configuration.Log(r.p.DbName, r.p.UserName, format, a...)
}

func (r ApiRequest) trace(format string, a ...any) {
	configuration.Trace(r.p.DbName, r.p.UserName, format, a...)
}

func (r ApiRequest) startTiming() time.Time {
	return configuration.StartTiming()
}

func (r ApiRequest) stopTiming(funcName string, start time.Time) {
	configuration.StopTiming(r.p.DbName, r.p.UserName, funcName, start)
}

func (r ApiRequest) logAndReturnError(format string, a ...any) error {
	return configuration.LogAndReturnError(r.p.DbName, r.p.UserName, format, a...)
}

func (r ApiRequest) execContext(query string, args ...any) error {
	_, err := r.db.ExecContext(r.ctx, query, args...)
	return err
}

func (r ApiRequest) queryContext(query string, args ...any) (*sql.Rows, error) {
	return r.db.QueryContext(r.ctx, query, args...)
}

func (r ApiRequest) beginTransaction() error {
	r.trace("beginTransaction()")
	if err := r.execContext(getBeginTransactionQuery()); err != nil {
		return r.logAndReturnError("Begin transaction error: %v.", err)
	}
	r.log("Begin transaction.")
	return nil
}

// function is available for mocking
var getBeginTransactionQuery = func() string {
	return "BEGIN"
}

func (r ApiRequest) commitTransaction() error {
	r.trace("commitTransaction()")
	if err := r.execContext(getCommitTransactionQuery()); err != nil {
		return r.logAndReturnError("Commit transaction error: %v.", err)
	}
	r.log("Commit transaction.")
	return nil
}

// function is available for mocking
var getCommitTransactionQuery = func() string {
	return "COMMIT"
}

func (r ApiRequest) rollbackTransaction() error {
	r.trace("rollbackTransaction()")
	if err := r.execContext(getRollbackTransactionQuery()); err != nil {
		return r.logAndReturnError("Rollback transaction error: %v.", err)
	}
	r.log("ROLLBACK transaction.")
	return nil
}

// function is available for mocking
var getRollbackTransactionQuery = func() string {
	return "ROLLBACK"
}

func createContextAndApiRequest(ct context.Context, p ApiParameters, uri string) (request ApiRequest, cancelFunc context.CancelFunc, error error) {
	ctx, cancel := configuration.GetContextWithTimeOut(ct, p.DbName)
	db, err := database.GetDbByName(p.DbName)
	r := ApiRequest{
		ctx:         ctx,
		p:           p,
		db:          db,
		ctxChan:     make(chan ApiResponse, 1),
		transaction: model.GetNewRowWithUuid(),
	}
	return r, cancel, err
}

func getHttpErrorCode(response ApiResponse) int {
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
	response := GetGridsRows(c.Request.Context(), c.Request.RequestURI, p)
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
