// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package core

import (
	"net/http"

	"d.lambert.fr/encoon/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/maps"
)

type user struct {
	Id        entity `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

var (
	users = make(map[string]user)
)

func loadUsers() {
	utils.Log("Loading users.")

	(&user{
		Id:        entity{Uuid: "c788a76d-4aa6-4073-8904-35a9b99a3289", Uri: "root", Version: 1, Enabled: true},
		Email:     "root@encoon.com",
		FirstName: "Root",
		LastName:  "Encoon"}).add()

	(&user{
		Id:        entity{Uuid: "bced42a2-6ddd-4023-ad40-0d46962b7872", Uri: "system", Version: 1, Enabled: true},
		Email:     "system@encoon.com",
		FirstName: "System",
		LastName:  "Encoon"}).add()

	(&user{
		Id:    entity{Uuid: "67b560b9-63ff-4fed-9b64-26c7f86e540c"},
		Email: "none@encoon.com"}).add()
}

func (v *user) add() {
	users[v.Id.Uuid] = *v
}

func GetUsersApi(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"users": maps.Values(users)})
}

func GetIndexHtml(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{"title": "εncooη"})
}

func GetUsersHtml(c *gin.Context) {
	c.HTML(http.StatusOK, "users.html", gin.H{"title": "Users", "users": users})
}

func GetUserByID(uuid string) (user, bool) {
	value, exists := users[uuid]
	return value, exists
}

func GetUserByIDApi(c *gin.Context) {
	uuid := c.Param("uuid")
	user, exists := GetUserByID(uuid)
	if exists {
		c.IndentedJSON(http.StatusOK, user)
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
	}
}

func PostUsersApi(c *gin.Context) {
	var newUser user
	if err := c.BindJSON(&newUser); err != nil {
		return
	}
	users[newUser.Id.Uuid] = newUser
	c.IndentedJSON(http.StatusCreated, newUser)
}
