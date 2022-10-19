// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoadData() {
	loadUsers()
}

func PingApi(c *gin.Context) {
	uuid := GetNewUUID()
	c.IndentedJSON(http.StatusOK, gin.H{"message": "ping", "uuid": uuid})
}
