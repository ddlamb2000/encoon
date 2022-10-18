// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package core

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/maps"
)

var users map[string]user

func loadUsers() {
	users = make(map[string]user)
	fmt.Fprintf(gin.DefaultWriter, "[εncooη] Loading users.\n")
	addUser(user{Id: entity{Uuid: "c788a76d-4aa6-4073-8904-35a9b99a3289", Uri: "root", Version: 1, Enabled: true}, Email: "root@encoon.com", FirstName: "Root", LastName: "Encoon"})
	addUser(user{Id: entity{Uuid: "bced42a2-6ddd-4023-ad40-0d46962b7872", Uri: "system", Version: 1, Enabled: true}, Email: "system@encoon.com", FirstName: "System", LastName: "Encoon"})
	addUser(user{Id: entity{Uuid: "67b560b9-63ff-4fed-9b64-26c7f86e540c"}, Email: "none@encoon.com"})
}

func addUser(user user) {
	users[user.Id.Uuid] = user
}

func GetUsersJson(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, maps.Values(users))
}

func GetUsersHtml(c *gin.Context) {
	c.HTML(http.StatusOK, "users.html", gin.H{"title": "Users", "users": users})
}

func GetUserByID(uuid string) (user, bool) {
	value, exists := users[uuid]
	return value, exists
}

func GetAlbumByIDJson(c *gin.Context) {
	uuid := c.Param("uuid")
	user, exists := GetUserByID(uuid)
	if exists {
		c.IndentedJSON(http.StatusOK, user)
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
	}
}

func PostUsers(c *gin.Context) {
	var newUser user
	if err := c.BindJSON(&newUser); err != nil {
		return
	}
	users[newUser.Id.Uuid] = newUser
	c.IndentedJSON(http.StatusCreated, newUser)
}
