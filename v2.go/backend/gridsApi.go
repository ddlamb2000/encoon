// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"d.lambert.fr/encoon/utils"
	"github.com/gin-gonic/gin"
)

func GetGridsApi(c *gin.Context) {
	time.Sleep(500 * time.Millisecond) ////// temporarisation that must be removed!
	dbName, gridUri, uuid, db := getDbForGridsApi(c)
	if db != nil {
		grid, err := getGridForGridsApi(db, gridUri)
		if err != nil {
			c.Abort()
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": err})
			return
		}
		getGridsApiAuthorized(c, dbName, db, gridUri, uuid, *grid)
	}
}

func getGridsApiAuthorized(c *gin.Context, dbName string, db *sql.DB, gridUri string, uuid string, grid Grid) {
	rows, err := getRowsForGridsApi(db, grid.Uuid, uuid)
	if err != nil {
		c.Abort()
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	defer rows.Close()
	rowSet, rowSetCount, err := getRowsetForGridsApi(dbName, gridUri, rows)
	if err != nil {
		c.Abort()
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"uuid": uuid, "grid": grid, "items": rowSet, "count": rowSetCount})
}

func getDbForGridsApi(c *gin.Context) (string, string, string, *sql.DB) {
	auth, exists := c.Get("authorized")
	if !exists || auth == false {
		c.Abort()
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Not authorized."})
		return "", "", "", nil
	}
	dbName := c.Param("dbName")
	gridUri := c.Param("gridUri")
	uuid := c.Param("uuid")
	if dbName == "" || gridUri == "" {
		c.Abort()
		c.IndentedJSON(http.StatusNotImplemented, gin.H{"error": "Missing parameter."})
		return "", "", "", nil
	}
	db := getDbByName(dbName)
	if db == nil {
		c.Abort()
		c.IndentedJSON(http.StatusNotImplemented, gin.H{"error": "Database isn't available."})
		return "", "", "", nil
	}
	return dbName, gridUri, uuid, db
}

func getGridForGridsApi(db *sql.DB, gridUri string) (*Grid, error) {
	grid := new(Grid)
	if err := db.QueryRow(
		"SELECT uuid, text01 FROM rows WHERE gridUuid = $1 AND uri = $2",
		utils.UuidGrids,
		gridUri).
		Scan(&grid.Uuid, &grid.Text01); err != nil {
		if err == sql.ErrNoRows {
			return grid, errors.New("GRID NOT FOUND")
		} else if grid.Uuid == "" {
			return grid, errors.New("GRID IDENTIFIER NOT FOUND")
		} else {
			return grid, fmt.Errorf("UNKNOWN ERROR WHEN RETRIEVING GRID DEFINITION: %v", err)
		}
	}
	return grid, nil
}

func getRowsQueryForGridsApi(uuid string) string {
	selectStr := "SELECT uuid, version, uri, text01, text02, text03, text04"
	fromStr := "FROM rows"
	whereStr := getRowsWhereQueryForGridsApi(uuid)
	return selectStr + " " + fromStr + " " + whereStr
}

func getRowsWhereQueryForGridsApi(uuid string) string {
	if uuid == "" {
		return "WHERE griduuid = $1"
	}
	return "WHERE uuid = $2 AND griduuid = $1"
}

func getRowsQueryParametersForGridsApi(gridUuid string, uuid string) []any {
	parameters := make([]any, 0)
	parameters = append(parameters, gridUuid)
	if uuid != "" {
		parameters = append(parameters, uuid)
	}
	return parameters
}

func getRowsForGridsApi(db *sql.DB, gridUuid string, uuid string) (*sql.Rows, error) {
	rows, err := db.Query(getRowsQueryForGridsApi(uuid), getRowsQueryParametersForGridsApi(gridUuid, uuid)...)
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
