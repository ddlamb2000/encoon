// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package middleware

import (
	"github.com/gin-gonic/gin"
)

func GetUsersApi(c *gin.Context) {
	// auth, exists := c.Get("authorized")
	// if !exists || auth == false {
	// 	c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
	// 	return
	// }
	// c.IndentedJSON(http.StatusOK, gin.H{"users": maps.Values(users)})
}

func GetUserByIDApi(c *gin.Context) {
	// auth, exists := c.Get("authorized")
	// if !exists || auth == false {
	// 	c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
	// 	return
	// }
	// uuid := c.Param("uuid")
	// user, exists := GetUserByID(uuid)
	// if exists {
	// 	c.IndentedJSON(http.StatusOK, gin.H{"users": user})
	// } else {
	// 	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
	// }
}

func PostUsersApi(c *gin.Context) {
	// auth, exists := c.Get("authorized")
	// if !exists || auth == false {
	// 	c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
	// 	return
	// }
	// var newUser user
	// if err := c.BindJSON(&newUser); err != nil {
	// 	return
	// }
	// users[newUser.entity.Uuid] = newUser
	// c.IndentedJSON(http.StatusCreated, newUser)
}
