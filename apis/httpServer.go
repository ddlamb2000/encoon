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

func getParameters(c *gin.Context) (htmlParameters, error) {
	dbName, gridUuid, uuid := c.Param("dbName"), c.Param("gridUuid"), c.Param("uuid")
	userUuid, userName := c.GetString("userUuid"), c.GetString("user")
	auth, exists := c.Get("authorized")
	p := htmlParameters{
		dbName:   dbName,
		userUuid: userUuid,
		userName: userName,
		gridUuid: gridUuid,
		uuid:     uuid,
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

type htmlParameters struct {
	dbName               string
	userUuid             string
	userName             string
	gridUuid             string
	uuid                 string
	filterColumnOwned    bool
	filterColumnName     string
	filterColumnGridUuid string
	filterColumnValue    string
}

type apiRequest struct {
	ctx         context.Context
	p           htmlParameters
	db          *sql.DB
	ctxChan     chan apiResponse
	transaction *model.Row
}

func (r apiRequest) log(format string, a ...any) {
	configuration.Log(r.p.dbName, r.p.userName, format, a...)
}

func (r apiRequest) trace(format string, a ...any) {
	configuration.Trace(r.p.dbName, r.p.userName, format, a...)
}

func (r apiRequest) startTiming() time.Time {
	return configuration.StartTiming()
}

func (r apiRequest) stopTiming(funcName string, start time.Time) {
	configuration.StopTiming(r.p.dbName, r.p.userName, funcName, start)
}

func (r apiRequest) logAndReturnError(format string, a ...any) error {
	return configuration.LogAndReturnError(r.p.dbName, r.p.userName, format, a...)
}

func (r apiRequest) execContext(query string, args ...any) error {
	_, err := r.db.ExecContext(r.ctx, query, args...)
	return err
}

func (r apiRequest) queryContext(query string, args ...any) (*sql.Rows, error) {
	return r.db.QueryContext(r.ctx, query, args...)
}

func (r apiRequest) beginTransaction() error {
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

func (r apiRequest) commitTransaction() error {
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

func (r apiRequest) rollbackTransaction() error {
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

type apiResponse struct {
	Grid                   *model.Grid         `json:"grid"`
	CountRows              int                 `json:"countRows"`
	Rows                   []model.Row         `json:"rows"`
	RowsAdded              []*model.Row        `json:"rowsAdded,omitempty"`
	RowsEdited             []*model.Row        `json:"rowsEdited,omitempty"`
	RowsDeleted            []*model.Row        `json:"rowsDeleted,omitempty"`
	ReferenceValuesAdded   []gridReferencePost `json:"referencedValuesAdded,omitempty"`
	ReferenceValuesRemoved []gridReferencePost `json:"referencedValuesRemoved,omitempty"`
	Err                    error               `json:"err,omitempty"`
	TimeOut                bool                `json:"timeOut,omitempty"`
	System                 bool                `json:"system,omitempty"`
	Forbidden              bool                `json:"forbidden,omitempty"`
	CanViewRows            bool                `json:"canViewRows"`
	CanEditRows            bool                `json:"canEditRows"`
	CanEditOwnedRows       bool                `json:"canEditOwnedRows"`
	CanAddRows             bool                `json:"canAddRows"`
}

func createContextAndApiRequest(ct context.Context, p htmlParameters, uri string) (request apiRequest, cancelFunc context.CancelFunc, error error) {
	ctx, cancel := configuration.GetContextWithTimeOut(ct, p.dbName)
	db, err := database.GetDbByName(p.dbName)
	r := apiRequest{
		ctx:         ctx,
		p:           p,
		db:          db,
		ctxChan:     make(chan apiResponse, 1),
		transaction: model.GetNewRowWithUuid(),
	}
	return r, cancel, err
}

func getHttpErrorCode(response apiResponse) int {
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
	response := getGridsRows(c.Request.Context(), c.Request.RequestURI, p)
	if response.Err != nil {
		c.Abort()
		c.JSON(getHttpErrorCode(response), gin.H{"error": response.Err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"response": response})
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
	ColumnUuid string `json:"columnUuid"`
	FromUuid   string `json:"fromUuid"`
	ToGridUuid string `json:"toGridUuid"`
	ToUuid     string `json:"uuid"`
	Owned      bool   `json:"owned"`
}

func PostGridsRowsApi(c *gin.Context) {
	p, err := getParameters(c)
	if err != nil {
		c.Abort()
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}
	var payload gridPost
	c.ShouldBindJSON(&payload)
	response := postGridsRows(c.Request.Context(), c.Request.RequestURI, p, payload)
	if response.Err != nil {
		c.Abort()
		c.JSON(getHttpErrorCode(response), gin.H{"error": response.Err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"response": response})
}
