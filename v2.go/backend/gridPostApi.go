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
	err := postGridsRows(dbName, userUuid, gridUri, payload.RowsAdded)
	if err != nil {
		c.Abort()
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "")
}

func postGridsRows(dbName string, userUuid string, gridUri string, rowsAdded []Row) error {
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
		row.Uuid = utils.GetNewUUID()
		err := postGridRow(ctx, dbName, db, userUuid, grid.Uuid, row)
		if err != nil {
			return err
		}
	}
	return nil
}

func postGridRow(ctx context.Context, dbName string, db *sql.DB, userUuid string, gridUuid string, row Row) error {
	insertStr := "INSERT INTO rows (uuid, " +
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
		"text04) " +
		"VALUES ($1, " +
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
	utils.Log("%v : %v", insertStr, *row.Uri)
	_, err := db.ExecContext(ctx, insertStr, row.Uuid, userUuid, gridUuid, row.Uri, row.Text01, row.Text02, row.Text03, row.Text04)
	if err != nil {
		utils.LogError("[%q] Insert row: %v", dbName, err)
	}
	return err
}
