// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import (
	"net/http"

	"d.lambert.fr/encoon/utils"
	"github.com/gin-gonic/gin"
)

type gridPost struct {
	RowsAdded  []Row `json:"rowsAdded"`
	RowsEdited []Row `json:"rowsEdited"`
}

func PostGridsRowsApi(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	auth, exists := c.Get("authorized")
	if !exists || auth == false {
		c.Abort()
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Not authorized."})
		return
	}
	dbName := c.Param("dbName")
	if dbName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing database parameter."})
		return
	}

	var payload gridPost
	c.BindJSON(&payload)
	utils.Log("%v", payload.RowsAdded)
	for _, row := range payload.RowsAdded {
		utils.Log("%v", *row.Uri)
	}
	// c.JSON(http.StatusOK, JWTtoken{tokenString})
}
