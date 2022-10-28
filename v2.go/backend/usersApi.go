// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/maps"
)

func GetUsersApi(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"users": maps.Values(users)})
}

func GetUserByIDApi(c *gin.Context) {
	uuid := c.Param("uuid")
	user, exists := GetUserByID(uuid)
	if exists {
		c.IndentedJSON(http.StatusOK, gin.H{"users": user})
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
	}
}

func PostUsersApi(c *gin.Context) {
	var newUser user
	if err := c.BindJSON(&newUser); err != nil {
		return
	}
	users[newUser.entity.Uuid] = newUser
	c.IndentedJSON(http.StatusCreated, newUser)
}
