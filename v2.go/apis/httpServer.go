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

func (r apiRequestParameters) queryContext(query string, args ...any) (*sql.Rows, error) {
	return r.db.QueryContext(r.ctx, query, args...)
}

func (r apiRequestParameters) beginTransaction() error {
	r.trace("beginTransaction()")
	if err := r.execContext(getBeginTransactionQuery()); err != nil {
		return r.logAndReturnError("Begin transaction error: %v.", err)
	}
	r.log("Begin transaction.")
	return nil
}

// function is available for mocking
var getBeginTransactionQuery = func() string { return "BEGIN" }

func (r apiRequestParameters) commitTransaction() error {
	r.trace("commitTransaction()")
	if err := r.execContext(getCommitTransactionQuery()); err != nil {
		return r.logAndReturnError("Commit transaction error: %v.", err)
	}
	r.log("Commit transaction.")
	return nil
}

// function is available for mocking
var getCommitTransactionQuery = func() string { return "COMMIT" }

func (r apiRequestParameters) rollbackTransaction() error {
	r.trace("rollbackTransaction()")
	if err := r.execContext(getRollbackTransactionQuery()); err != nil {
		return r.logAndReturnError("Rollback transaction error: %v.", err)
	}
	r.log("ROLLBACK transaction.")
	return nil
}

// function is available for mocking
var getRollbackTransactionQuery = func() string { return "ROLLBACK" }

type apiResponse struct {
	grid     *model.Grid
	rows     []model.Row
	rowCount int
	err      error
	timeOut  bool
	system   bool
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

func getHttpErrorCode(response apiResponse) int {
	if response.system {
		return http.StatusInternalServerError
	} else if response.timeOut {
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
	grid, rowSet, rowSetCount, response := getGridsRows(c.Request.Context(), dbName, gridUuid, uuid, userUuid, userName)
	if response.err != nil {
		c.Abort()
		c.JSON(getHttpErrorCode(response), gin.H{"error": response.err.Error()})
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
	grid, rowSet, rowSetCount, response := postGridsRows(c.Request.Context(), dbName, userUuid, userName, gridUuid, uuid, payload)
	if response.err != nil {
		c.Abort()
		c.JSON(getHttpErrorCode(response), gin.H{"error": response.err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"grid": grid, "uuid": uuid, "rows": rowSet, "countRows": rowSetCount})
}
