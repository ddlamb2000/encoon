// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"d.lambert.fr/encoon/utils"
	"github.com/gin-gonic/gin"
)

func GetGridsApi(c *gin.Context) {
	time.Sleep(500 * time.Millisecond) ////// temporarisation that must be removed!
	auth, exists := c.Get("authorized")
	if !exists || auth == false {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Not authorized."})
		return
	}
	dbName := c.Param("dbName")
	gridUri := c.Param("gridUri")
	if dbName == "" || gridUri == "" {
		c.IndentedJSON(http.StatusNotImplemented, gin.H{"error": "Missing parameter."})
		return
	}
	db := getDbByName(dbName)
	if db == nil {
		c.IndentedJSON(http.StatusNotImplemented, gin.H{"error": "Database isn't available."})
		return
	}
	getGridsApiAuthorized(c, dbName, db, gridUri)
}

func getGridsApiAuthorized(c *gin.Context, dbName string, db *sql.DB, gridUri string) {
	uuid := c.Param("uuid")
	var err error
	var grid Grid
	if err = db.QueryRow(
		"SELECT uuid, text01 FROM rows WHERE gridUuid = $1 AND uri = $2",
		utils.UuidGrids,
		gridUri).
		Scan(&grid.Uuid, &grid.Text01); err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound,
				gin.H{"error": "Grid not found."})
			return
		} else if grid.Uuid == "" {
			c.IndentedJSON(http.StatusNotFound,
				gin.H{"error": "Grid identifier not found."})
			return
		} else {
			c.IndentedJSON(http.StatusInternalServerError,
				gin.H{"error": fmt.Sprintf("Unknown error when retrieving grid definition: %v.", err)})
			return
		}
	}

	var rows *sql.Rows
	if uuid != "" {
		rows, err = db.Query(
			"SELECT uuid, version, uri, text01, text02, text03, text04 FROM rows WHERE uuid = $1 AND griduuid = $2",
			uuid,
			grid.Uuid)

	} else {
		rows, err = db.Query(
			"SELECT uuid, version, uri, text01, text02, text03, text04 FROM rows WHERE griduuid = $1",
			grid.Uuid)
	}

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("Error when querying rows: %v.", err)})
		return
	}
	defer rows.Close()

	var rowSet = make([]Row, 0)
	var rowSetCount = 0
	for rows.Next() {
		var row Row
		err = rows.Scan(&row.Uuid, &row.Version, &row.Uri, &row.Text01, &row.Text02, &row.Text03, &row.Text04)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError,
				gin.H{"error": fmt.Sprintf("Unknown error when scanning rows: %v.", err)})
			return
		}
		row.Path = fmt.Sprintf("/%s/%s/%s", dbName, gridUri, row.Uuid)
		rowSet = append(rowSet, row)
		rowSetCount += 1
	}
	err = rows.Err()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("Error when scanning rows: %v.", err)})
		return
	}

	c.IndentedJSON(http.StatusOK,
		gin.H{
			"uuid":  uuid,
			"grid":  grid,
			"items": rowSet,
			"count": rowSetCount,
		})
}
