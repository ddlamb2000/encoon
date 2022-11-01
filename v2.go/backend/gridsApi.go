// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"d.lambert.fr/encoon/utils"
	"github.com/gin-gonic/gin"
)

func GetGridsRowsApi(c *gin.Context) {
	auth, exists := c.Get("authorized")
	if !exists || auth == false {
		c.Abort()
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Not authorized."})
		return
	}
	dbName := c.Param("dbName")
	gridUri := c.Param("gridUri")
	uuid := c.Param("uuid")
	grid, rowSet, rowSetCount, err := getGridsRows(dbName, gridUri, uuid)
	if err != nil {
		c.Abort()
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}
	c.IndentedJSON(http.StatusOK, gin.H{"uuid": uuid, "grid": grid, "items": rowSet, "count": rowSetCount})
}

func getGridsRows(dbName string, gridUri string, uuid string) (*Grid, []Row, int, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(utils.Configuration.DbTimeOut)*time.Second)
	defer cancel()

	db, err := getDbForGridsApi(dbName, gridUri)
	if err != nil {
		return nil, nil, 0, err
	}
	grid, err := getGridForGridsApi(ctx, db, gridUri)
	if err != nil {
		return nil, nil, 0, err
	}
	rows, err := getRowsForGridsApi(ctx, db, grid.Uuid, uuid)
	if err != nil {
		return nil, nil, 0, err
	}
	defer rows.Close()
	rowSet, rowSetCount, err := getRowsetForGridsApi(dbName, gridUri, rows)

	return grid, rowSet, rowSetCount, err
}

func getDbForGridsApi(dbName string, gridUri string) (*sql.DB, error) {
	if dbName == "" || gridUri == "" {
		return nil, fmt.Errorf("MISSING PARAMETER")
	}
	db := getDbByName(dbName)
	if db == nil {
		return nil, fmt.Errorf("DATABASE NOT AVAILABLE")
	}
	return db, nil
}

func getGridForGridsApi(ctx context.Context, db *sql.DB, gridUri string) (*Grid, error) {
	grid := new(Grid)
	if err := db.QueryRowContext(
		ctx,
		"SELECT uuid, text01 FROM rows WHERE gridUuid = $1 AND uri = $2",
		utils.UuidGrids,
		gridUri).
		Scan(
			&grid.Uuid,
			&grid.Text01); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("GRID NOT FOUND")
		} else if grid.Uuid == "" {
			return nil, errors.New("GRID IDENTIFIER NOT FOUND")
		} else {
			return nil, fmt.Errorf("UNKNOWN ERROR WHEN RETRIEVING GRID DEFINITION: %v", err)
		}
	}
	return grid, nil
}

func getRowsQueryForGridsApi(uuid string) string {
	selectStr := " SELECT uuid, version, uri, text01, text02, text03, text04 "
	fromStr := " FROM rows "
	whereStr := getRowsWhereQueryForGridsApi(uuid)
	return selectStr + fromStr + whereStr
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

func getRowsForGridsApi(ctx context.Context, db *sql.DB, gridUuid string, uuid string) (*sql.Rows, error) {
	rows, err := db.QueryContext(ctx,
		getRowsQueryForGridsApi(uuid),
		getRowsQueryParametersForGridsApi(gridUuid, uuid)...)
	if err != nil {
		return nil, fmt.Errorf("ERROR WHEN QUERYING ROWS: %v", err)
	}
	return rows, nil
}

func getRowsQueryOutputForGridsApi(row *Row) []any {
	output := make([]any, 0)
	output = append(output, &row.Uuid)
	output = append(output, &row.Version)
	output = append(output, &row.Uri)
	output = append(output, &row.Text01)
	output = append(output, &row.Text02)
	output = append(output, &row.Text03)
	output = append(output, &row.Text04)
	return output

}

func getRowsetForGridsApi(dbName string, gridUri string, rows *sql.Rows) ([]Row, int, error) {
	var rowSet = make([]Row, 0)
	var rowSetCount = 0
	for rows.Next() {
		var row = new(Row)
		if err := rows.Scan(getRowsQueryOutputForGridsApi(row)...); err != nil {
			return nil, 0, fmt.Errorf("UNKNOWN ERROR WHEN SCANNING ROWS: %v", err)
		}
		row.SetPath(dbName, gridUri)
		rowSet = append(rowSet, *row)
		rowSetCount += 1
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("ERROR WHEN SCANNNIG ROWS: %v", err)
	}
	return rowSet, rowSetCount, nil
}
