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
	RowsAdded   []Row `json:"rowsAdded"`
	RowsEdited  []Row `json:"rowsEdited"`
	RowsDeleted []Row `json:"rowsDeleted"`
}

func PostGridsRowsApi(c *gin.Context) {
	userUuid, user, err := getUserUui(c)
	if err != nil {
		c.Abort()
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}
	dbName := c.Param("dbName")
	gridUri := c.Param("gridUri")
	logUri(c, dbName, user)
	testSleep(dbName)
	var payload gridPost
	c.ShouldBindJSON(&payload)
	err = postGridsRows(dbName, userUuid, user, gridUri, payload.RowsAdded, payload.RowsEdited, payload.RowsDeleted)
	if err != nil {
		c.Abort()
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	grid, rowSet, rowSetCount, err := getGridsRows(dbName, gridUri, "", user)
	if err != nil {
		c.Abort()
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"grid": grid, "rows": rowSet, "countRows": rowSetCount})
}

func postGridsRows(dbName string, userUuid string, user string, gridUri string, rowsAdded []Row, rowsEdited []Row, rowsDeleted []Row) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(utils.Configuration.TimeOutThreshold)*time.Millisecond)
	defer cancel()

	db, err := getDbForGridsApi(dbName, user)
	if err != nil {
		return err
	}
	grid, err := getGridForGridsApi(ctx, db, dbName, user, gridUri)
	if err != nil {
		return err
	}
	if err := beginTransaction(ctx, dbName, db, userUuid, user); err != nil {
		return err
	}
	for _, row := range rowsAdded {
		err := postInsertGridRow(ctx, dbName, db, userUuid, user, gridUri, grid.Uuid, row)
		if err != nil {
			_ = rollbackTransaction(ctx, dbName, db, userUuid, user)
			return err
		}
	}
	for _, row := range rowsEdited {
		err := postUpdateGridRow(ctx, dbName, db, userUuid, user, gridUri, grid.Uuid, row)
		if err != nil {
			_ = rollbackTransaction(ctx, dbName, db, userUuid, user)
			return err
		}
	}
	for _, row := range rowsDeleted {
		err := postDeleteGridRow(ctx, dbName, db, userUuid, user, gridUri, grid.Uuid, row)
		if err != nil {
			_ = rollbackTransaction(ctx, dbName, db, userUuid, user)
			return err
		}
	}
	if err := commitTransaction(ctx, dbName, db, userUuid, user); err != nil {
		return err
	}
	return nil
}

func postInsertGridRow(ctx context.Context, dbName string, db *sql.DB, userUuid string, user string, gridUri string, gridUuid string, row Row) error {
	row.Uuid = utils.GetNewUUID()
	insertStatement := getInsertStatementForGridsApi()
	insertValues := getInsertValuesForGridsApi(userUuid, gridUuid, row)
	_, err := db.ExecContext(ctx, insertStatement, insertValues...)
	if err != nil {
		return utils.LogAndReturnError("[%s] [%s] Insert row error on %q: %v.", dbName, user, insertStatement, err)
	}
	utils.Log("[%s] [%s] Row [%s] inserted into %q.", dbName, user, row, gridUri)
	return err
}

func getInsertStatementForGridsApi() string {
	insertStr := "INSERT INTO rows (uuid, " +
		"version, " +
		"created, " +
		"updated, " +
		"createdBy, " +
		"updatedBy, " +
		"enabled, " +
		"gridUuid, " +
		"text01, " +
		"text02, " +
		"text03, " +
		"text04, " +
		"int01, " +
		"int02, " +
		"int03, " +
		"int04) "
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
		"$8, " +
		"$9, " +
		"$10, " +
		"$11)"
	return insertStr + valueStr
}

func getInsertValuesForGridsApi(userUuid string, gridUuid string, row Row) []any {
	values := make([]any, 0)
	values = append(values, row.Uuid)
	values = append(values, userUuid)
	values = append(values, gridUuid)
	values = append(values, row.Text01)
	values = append(values, row.Text02)
	values = append(values, row.Text03)
	values = append(values, row.Text04)
	values = append(values, row.Int01)
	values = append(values, row.Int02)
	values = append(values, row.Int03)
	values = append(values, row.Int04)
	return values
}

func postUpdateGridRow(ctx context.Context, dbName string, db *sql.DB, userUuid string, user string, gridUri string, gridUuid string, row Row) error {
	updateStatement := getUpdateStatementForGridsApi()
	updateValues := getUpdateValuesForGridsApi(userUuid, gridUuid, row)
	_, err := db.ExecContext(ctx, updateStatement, updateValues...)
	if err != nil {
		return utils.LogAndReturnError("[%s] [%s] Update row error on %q: %v.", dbName, user, updateStatement, err)
	}
	utils.Log("[%s] [%s] Row [%s] updated in %q.", dbName, user, row, gridUri)
	return err
}

func getUpdateStatementForGridsApi() string {
	updateStr := "UPDATE rows SET " +
		"version = version + 1, " +
		"updated = NOW(), " +
		"updatedBy = $3, " +
		"text01 = $4, " +
		"text02 = $5, " +
		"text03 = $6, " +
		"text04 = $7, " +
		"int01 = $8, " +
		"int02 = $9, " +
		"int03 = $10, " +
		"int04 = $11"
	whereStr := " WHERE uuid = $1 and gridUuid = $2"
	return updateStr + whereStr
}

func getUpdateValuesForGridsApi(userUuid string, gridUuid string, row Row) []any {
	values := make([]any, 0)
	values = append(values, row.Uuid)
	values = append(values, gridUuid)
	values = append(values, userUuid)
	values = append(values, row.Text01)
	values = append(values, row.Text02)
	values = append(values, row.Text03)
	values = append(values, row.Text04)
	values = append(values, row.Int01)
	values = append(values, row.Int02)
	values = append(values, row.Int03)
	values = append(values, row.Int04)
	return values
}

func postDeleteGridRow(ctx context.Context, dbName string, db *sql.DB, userUuid string, user string, gridUri string, gridUuid string, row Row) error {
	_, err := db.ExecContext(ctx, "DELETE FROM rows WHERE uuid = $1 and gridUuid = $2", row.Uuid, gridUuid)
	if err != nil {
		return utils.LogAndReturnError("[%s] [%s] Delete row error: %v.", dbName, user, err)
	}
	utils.Log("[%s] [%s] Row [%s] deleted in %q.", dbName, user, row, gridUri)
	return err
}
