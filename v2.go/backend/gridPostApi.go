// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"d.lambert.fr/encoon/utils"
	"github.com/gin-gonic/gin"
)

type gridPost struct {
	RowsAdded  []Row `json:"rowsAdded"`
	RowsEdited []Row `json:"rowsEdited"`
}

func PostGridsRowsApi(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	if !isAuthorized(c) {
		return
	}
	userUuid := getUserUui(c)
	if userUuid == "" {
		return
	}
	dbName := c.Param("dbName")
	gridUri := c.Param("gridUri")
	var payload gridPost
	c.BindJSON(&payload)
	err := postGridsRows(dbName, userUuid, gridUri, payload.RowsAdded, payload.RowsEdited)
	if err != nil {
		c.Abort()
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "")
}

func postGridsRows(dbName string, userUuid string, gridUri string, rowsAdded []Row, rowsEdited []Row) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(utils.Configuration.DbTimeOut)*time.Second)
	defer cancel()

	db, err := getDbForGridsApi(dbName, gridUri)
	if err != nil {
		return err
	}
	grid, err := getGridForGridsApi(ctx, db, gridUri)
	if err != nil {
		return err
	}

	for _, row := range rowsAdded {
		err := postInsertGridRow(ctx, dbName, db, userUuid, grid.Uuid, row)
		if err != nil {
			return err
		}
	}
	for _, row := range rowsEdited {
		err := postUpdateGridRow(ctx, dbName, db, userUuid, grid.Uuid, row)
		if err != nil {
			return err
		}
	}
	return nil
}

func postInsertGridRow(ctx context.Context, dbName string, db *sql.DB, userUuid string, gridUuid string, row Row) error {
	row.Uuid = utils.GetNewUUID()
	insertStatement := getInsertStatementForGridsApi()
	insertValues := getInsertValuesForGridsApi(userUuid, gridUuid, row)
	_, err := db.ExecContext(ctx, insertStatement, insertValues...)
	if err != nil {
		utils.LogError("[%q] Insert row error on %q: %v", dbName, insertStatement, err)
	}
	return err
}

func getInsertStatementForGridsApi() string {
	insertStr := " INSERT INTO rows (uuid, " +
		"version, " +
		"created, " +
		"updated, " +
		"createdBy, " +
		"updatedBy, " +
		"enabled, " +
		"gridUuid, " +
		"uri, " +
		"text01, " +
		"text02, " +
		"text03, " +
		"text04) "
	valueStr := " VALUES ($1, " +
		"1, " +
		"NOW(), " +
		"NOW(), " +
		"$2, " +
		"$2, " +
		"true, " +
		"$3, " +
		"$4, " +
		"$5, " +
		"$6, " +
		"$7, " +
		"$8)"
	return insertStr + valueStr
}

func getInsertValuesForGridsApi(userUuid string, gridUuid string, row Row) []any {
	values := make([]any, 0)
	values = append(values, row.Uuid)
	values = append(values, userUuid)
	values = append(values, gridUuid)
	values = append(values, row.Uri)
	values = append(values, row.Text01)
	values = append(values, row.Text02)
	values = append(values, row.Text03)
	values = append(values, row.Text04)
	return values
}

func postUpdateGridRow(ctx context.Context, dbName string, db *sql.DB, userUuid string, gridUuid string, row Row) error {
	updateStatement := getUpdateStatementForGridsApi()
	updateValues := getUpdateValuesForGridsApi(userUuid, gridUuid, row)
	_, err := db.ExecContext(ctx, updateStatement, updateValues...)
	if err != nil {
		utils.LogError("[%q] Update row error on %q: %v", dbName, updateStatement, err)
	}
	return err
}

func getUpdateStatementForGridsApi() string {
	updateStr := " UPDATE rows SET " +
		"version = version + 1, " +
		"updated = NOW(), " +
		"updatedBy = $3, " +
		"uri = $4, " +
		"text01 = $5, " +
		"text02 = $6, " +
		"text03 = $7, " +
		"text04 = $8 "
	whereStr := " WHERE uuid = $1 and gridUuid = $2 "
	return updateStr + whereStr
}

func getUpdateValuesForGridsApi(userUuid string, gridUuid string, row Row) []any {
	values := make([]any, 0)
	values = append(values, row.Uuid)
	values = append(values, gridUuid)
	values = append(values, userUuid)
	values = append(values, row.Uri)
	values = append(values, row.Text01)
	values = append(values, row.Text02)
	values = append(values, row.Text03)
	values = append(values, row.Text04)
	return values
}
