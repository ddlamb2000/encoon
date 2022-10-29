// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package middleware

import (
	"database/sql"
	"net/http"

	"d.lambert.fr/encoon/backend"
	"d.lambert.fr/encoon/utils"
	"github.com/gin-gonic/gin"
)

func GetGridsApi(c *gin.Context) {
	auth, exists := c.Get("authorized")
	if !exists || auth == false {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
		return
	}
	dbName := c.Param("dbName")
	gridUri := c.Param("gridUri")
	uuid := c.Param("uuid")
	if dbName == "" || gridUri == "" {
		c.IndentedJSON(http.StatusNotImplemented, gin.H{"message": "Missing parameter."})
		return
	}
	db := getDbByName(dbName)
	if db == nil {
		c.IndentedJSON(http.StatusNotImplemented, gin.H{"message": "Database isn't available."})
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
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Grid not found."})
			return
		} else if gridUuid == "" {
			utils.Log("[%q] Grid not found.", dbName)
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Grid not found."})
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

	var rowSet = make([]backend.Row, 0)
	for rows.Next() {
		var row backend.Row
		err = rows.Scan(&row.Uuid, &row.Version, &row.Uri, &row.Text01, &row.Text02, &row.Text03, &row.Text04)
		if err != nil {
			utils.Log("[%q] Unknown error when scanning rows: %v.", dbName, err)
			c.IndentedJSON(http.StatusInternalServerError, "")
			return
		}
		rowSet = append(rowSet, row)
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
		})
}
