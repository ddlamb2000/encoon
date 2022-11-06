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

func GetGridsRowsApi(c *gin.Context) {
	_, user, err := getUserUui(c)
	if err != nil {
		c.Abort()
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	dbName := c.Param("dbName")
	gridUri := c.Param("gridUri")
	uuid := c.Param("uuid")
	logUri(c, dbName, user)
	grid, rowSet, rowSetCount, err := getGridsRows(dbName, gridUri, uuid, user)
	if err != nil {
		c.Abort()
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	grid.SetPath(dbName, gridUri)
	c.IndentedJSON(http.StatusOK, gin.H{"uuid": uuid, "grid": grid, "rows": rowSet, "countRows": rowSetCount})
}

func getUserUui(c *gin.Context) (string, string, error) {
	auth, exists := c.Get("authorized")
	if !exists || auth == false {
		return "", "", utils.LogAndReturnError("Not authorized for %v.", c.Request.URL)
	}
	userUuid, user := c.GetString("userUuid"), c.GetString("user")
	if userUuid == "" || user == "" {
		return "", "", utils.LogAndReturnError("User not authorized for %v.", c.Request.URL)
	}
	return userUuid, user, nil
}

func getGridsRows(dbName string, gridUri string, uuid string, user string) (*Grid, []Row, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(utils.Configuration.DbTimeOut)*time.Second)
	defer cancel()

	db, err := getDbForGridsApi(dbName, user)
	if err != nil {
		return nil, nil, 0, err
	}
	grid, err := getGridForGridsApi(ctx, db, dbName, user, gridUri)
	if err != nil {
		return nil, nil, 0, err
	}
	rows, err := getRowsForGridsApi(ctx, db, dbName, user, grid.Uuid, uuid)
	if err != nil {
		return nil, nil, 0, err
	}
	defer rows.Close()
	rowSet, rowSetCount, err := getRowSetForGridsApi(dbName, user, gridUri, rows)

	return grid, rowSet, rowSetCount, err
}

func getDbForGridsApi(dbName string, user string) (*sql.DB, error) {
	if dbName == "" {
		return nil, utils.LogAndReturnError("[%s] Missing database name parameter.", user)
	}
	db := getDbByName(dbName)
	if db == nil {
		return nil, utils.LogAndReturnError("[%s] [%s] Database not available.", dbName, user)
	}
	return db, nil
}

func getGridForGridsApi(ctx context.Context,
	db *sql.DB,
	dbName string,
	user string,
	gridUri string) (*Grid, error) {
	selectGridStatement := "SELECT uuid, text01 FROM rows WHERE gridUuid = $1 AND text01 = $2"
	grid := new(Grid)
	if err := db.QueryRowContext(ctx, selectGridStatement, utils.UuidGrids, gridUri).
		Scan(&grid.Uuid, &grid.Text01); err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.LogAndReturnError("[%s] [%s] Grid %q not found.", dbName, user, gridUri)
		} else {
			return nil, utils.LogAndReturnError("[%s] [%s] Error when retrieving grid definition %q: %v.", dbName, user, gridUri, err)
		}
	}
	utils.Log("[%s] [%s] Got grid %q: [%s].", dbName, user, gridUri, grid)
	return grid, nil
}

func getRowsForGridsApi(ctx context.Context, db *sql.DB, dbName string, user string, gridUuid string, uuid string) (*sql.Rows, error) {
	rows, err := db.QueryContext(ctx,
		getRowsQueryForGridsApi(uuid),
		getRowsQueryParametersForGridsApi(gridUuid, uuid)...)
	if err != nil {
		return nil, utils.LogAndReturnError("[%s] [%s] Error when querying rows: %v.", dbName, user, err)
	}
	return rows, nil
}

func getRowsQueryForGridsApi(uuid string) string {
	selectStr := " SELECT uuid, version, text01, text02, text03, text04 "
	fromStr := " FROM rows "
	whereStr := getRowsWhereQueryForGridsApi(uuid)
	orderByStr := " ORDER BY text01, text02, text03, text04 "
	return selectStr + fromStr + whereStr + orderByStr
}

func getRowsWhereQueryForGridsApi(uuid string) string {
	if uuid == "" {
		return " WHERE griduuid = $1 "
	}
	return " WHERE uuid = $2 AND griduuid = $1 "
}

func getRowsQueryParametersForGridsApi(gridUuid string, uuid string) []any {
	parameters := make([]any, 0)
	parameters = append(parameters, gridUuid)
	if uuid != "" {
		parameters = append(parameters, uuid)
	}
	return parameters
}

func getRowSetForGridsApi(dbName string, user string, gridUri string, rows *sql.Rows) ([]Row, int, error) {
	var rowSet = make([]Row, 0)
	var rowSetCount = 0
	for rows.Next() {
		var row = new(Row)
		if err := rows.Scan(getRowsQueryOutputForGridsApi(row)...); err != nil {
			return nil, 0, utils.LogAndReturnError("[%s] [%s] Error when scanning rows for %q: %v.", dbName, user, gridUri, err)
		}
		row.SetPath(dbName, gridUri)
		rowSet = append(rowSet, *row)
		rowSetCount += 1
	}
	if err := rows.Err(); err != nil {
		return nil, 0, utils.LogAndReturnError("[%s] [%s] Error when scanning rows for %q: %v.", dbName, user, gridUri, err)
	}
	utils.Log("[%s] [%s] Got rows from %q.", dbName, user, gridUri)
	return rowSet, rowSetCount, nil
}

func getRowsQueryOutputForGridsApi(row *Row) []any {
	output := make([]any, 0)
	output = append(output, &row.Uuid)
	output = append(output, &row.Version)
	output = append(output, &row.Text01)
	output = append(output, &row.Text02)
	output = append(output, &row.Text03)
	output = append(output, &row.Text04)
	return output
}
