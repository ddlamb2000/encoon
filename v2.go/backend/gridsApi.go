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

	time.Sleep(1000 * time.Millisecond) ////// temporarisation that must be removed!

	auth, exists := c.Get("authorized")
	if !exists || auth == false {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Not authorized."})
		return
	}
	dbName := c.Param("dbName")
	gridUri := c.Param("gridUri")
	uuid := c.Param("uuid")
	if dbName == "" || gridUri == "" {
		c.IndentedJSON(http.StatusNotImplemented, gin.H{"error": "Missing parameter."})
		return
	}
	db := getDbByName(dbName)
	if db == nil {
		c.IndentedJSON(http.StatusNotImplemented, gin.H{"error": "Database isn't available."})
		return
	}

	var gridUuid string
	if err := db.QueryRow(
		"SELECT uuid FROM rows WHERE gridUuid = $1 AND uri = $2",
		utils.UuidGrids,
		gridUri).
		Scan(&gridUuid); err != nil {
		if err == sql.ErrNoRows {
			utils.Log("[%q] Grid not found: %v.", dbName, err)
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Grid not found."})
			return
		} else if gridUuid == "" {
			utils.Log("[%q] Grid not found.", dbName)
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Grid not found."})
			return
		} else {
			utils.Log("[%q] Unknown error: %v.", dbName, err)
			c.IndentedJSON(http.StatusInternalServerError, "")
			return
		}
	}

	var rows *sql.Rows
	var err error
	if uuid != "" {
		rows, err = db.Query("SELECT uuid, version, uri, text01, text02, text03, text04 FROM rows WHERE uuid = $1 AND griduuid = $2", uuid, gridUuid)

	} else {
		rows, err = db.Query("SELECT uuid, version, uri, text01, text02, text03, text04 FROM rows WHERE griduuid = $1", gridUuid)
	}

	if err != nil {
		utils.Log("[%q] Error when querying rows: %v.", dbName, err)
		c.IndentedJSON(http.StatusInternalServerError, "")
		return
	}
	defer rows.Close()

	var rowSet = make([]Row, 0)
	var rowSetCount = 0
	for rows.Next() {
		var row Row
		err = rows.Scan(&row.Uuid, &row.Version, &row.Uri, &row.Text01, &row.Text02, &row.Text03, &row.Text04)
		if err != nil {
			utils.Log("[%q] Unknown error when scanning rows: %v.", dbName, err)
			c.IndentedJSON(http.StatusInternalServerError, "")
			return
		}
		row.Path = fmt.Sprintf("/%s/%s/%s", dbName, gridUri, row.Uuid)
		rowSet = append(rowSet, row)
		rowSetCount += 1
	}
	err = rows.Err()
	if err != nil {
		utils.Log("[%q] Error when scanning rows: %v.", dbName, err)
		c.IndentedJSON(http.StatusInternalServerError, "")
		return
	}

	c.IndentedJSON(http.StatusOK,
		gin.H{
			"gridUri":  gridUri,
			"gridUuid": gridUuid,
			"uuid":     uuid,
			"items":    rowSet,
			"count":    rowSetCount,
		})
}
