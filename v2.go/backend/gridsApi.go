// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/maps"
)

func GetGridsApi(c *gin.Context) {
	auth, exists := c.Get("authorized")
	if !exists || auth == false {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"users": maps.Values(users)})
}
