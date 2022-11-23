// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"context"
	"database/sql"
	"net/http"

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
	}
}

func getParameters(c *gin.Context) (string, string, string, string, string, error) {
	dbName, gridUuid, uuid := c.Param("dbName"), c.Param("gridUuid"), c.Param("uuid")
	userUuid, userName := c.GetString("userUuid"), c.GetString("user")
	auth, exists := c.Get("authorized")
	if dbName == "" || gridUuid == "" || len(userUuid) != len(model.UuidUsers) || userName == "" || !exists || auth == false {
		return "", "", "", "", "", configuration.LogAndReturnError(dbName, userName, "User not authorized.")
	}
	return dbName, userUuid, userName, gridUuid, uuid, nil
}

type apiRequestParameters struct {
	ctx      context.Context
	dbName   string
	userName string
	userUuid string
	db       *sql.DB
	ctxChan  chan apiResponse
}

func (r apiRequestParameters) log(format string, a ...any) {
	configuration.Log(r.dbName, r.userName, format, a...)
}

func (r apiRequestParameters) trace(format string, a ...any) {
	configuration.Trace(r.dbName, r.userName, format, a...)
}

func (r apiRequestParameters) logAndReturnError(format string, a ...any) error {
	return configuration.LogAndReturnError(r.dbName, r.userName, format, a...)
}

func (r apiRequestParameters) execContext(query string, args ...any) error {
	_, err := r.db.ExecContext(r.ctx, query, args...)
	return err
}

func (r apiRequestParameters) beginTransaction() error {
	r.trace("beginTransaction()")
	if err := r.execContext("BEGIN"); err != nil {
		return r.logAndReturnError("Begin transaction error: %v.", err)
	}
	r.log("Begin transaction.")
	return nil
}

func (r apiRequestParameters) commitTransaction() error {
	r.trace("commitTransaction()")
	if err := r.execContext("COMMIT"); err != nil {
		return r.logAndReturnError("Commit transaction error: %v.", err)
	}
	r.log("Commit transaction.")
	return nil
}

func (r apiRequestParameters) rollbackTransaction() error {
	r.trace("rollbackTransaction()")
	if err := r.execContext("ROLLBACK"); err != nil {
		return r.logAndReturnError("Rollback transaction error: %v.", err)
	}
	r.log("ROLLBACK transaction.")
	return nil
}

type apiResponse struct {
	grid     *model.Grid
	rows     []model.Row
	rowCount int
	err      error
}

func createContextAndApiRequestParameters(ct context.Context, dbName, userUuid, user string) (apiRequestParameters, context.CancelFunc, error) {
	ctx, cancel := configuration.GetContextWithTimeOut(ct, dbName)
	db, err := database.GetDbByName(dbName)
	r := apiRequestParameters{
		ctx:      ctx,
		dbName:   dbName,
		userName: user,
		userUuid: userUuid,
		db:       db,
		ctxChan:  make(chan apiResponse, 1),
	}
	return r, cancel, err
}

func getHttpErrorCode(timeOut bool) int {
	if timeOut {
		return http.StatusRequestTimeout
	} else {
		return http.StatusNotFound
	}
}

func GetGridsRowsApi(c *gin.Context) {
	dbName, userUuid, userName, gridUuid, uuid, err := getParameters(c)
	if err != nil {
		c.Abort()
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	grid, rowSet, rowSetCount, timeOut, err := getGridsRows(c.Request.Context(), dbName, gridUuid, uuid, userUuid, userName)
	if err != nil {
		c.Abort()
		c.JSON(getHttpErrorCode(timeOut), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"grid": grid, "uuid": uuid, "rows": rowSet, "countRows": rowSetCount})
}

type gridPost struct {
	RowsAdded              []*model.Row        `json:"rowsAdded"`
	RowsEdited             []*model.Row        `json:"rowsEdited"`
	RowsDeleted            []*model.Row        `json:"rowsDeleted"`
	ReferenceValuesAdded   []gridReferencePost `json:"referencedValuesAdded"`
	ReferenceValuesRemoved []gridReferencePost `json:"referencedValuesRemoved"`
}

type gridReferencePost struct {
	ColumnName string `json:"columnName"`
	FromUuid   string `json:"fromUuid"`
	ToGridUuid string `json:"toGridUuid"`
	ToUuid     string `json:"uuid"`
}

func PostGridsRowsApi(c *gin.Context) {
	dbName, userUuid, userName, gridUuid, uuid, err := getParameters(c)
	if err != nil {
		c.Abort()
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}
	var payload gridPost
	c.ShouldBindJSON(&payload)
	grid, rowSet, rowSetCount, timeOut, err := postGridsRows(c.Request.Context(), dbName, userUuid, userName, gridUuid, uuid, payload)
	if err != nil {
		c.Abort()
		c.JSON(getHttpErrorCode(timeOut), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"grid": grid, "uuid": uuid, "rows": rowSet, "countRows": rowSetCount})
}
