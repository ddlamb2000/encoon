// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package core

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/maps"
)

type user struct {
	ID        entity `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

var users map[string]user

var NoUser user = user{}

func LoadUsers() {
	fmt.Fprintf(gin.DefaultWriter, "[encoon] Loading users.\n")
	users = make(map[string]user)
	users["c788a76d-4aa6-4073-8904-35a9b99a3289"] = user{ID: entity{Uuid: "c788a76d-4aa6-4073-8904-35a9b99a3289", Uri: "root", Version: 1, Enabled: true}, Email: "root@encoon.com", FirstName: "Root", LastName: "Encoon"}
}

func GetUsersJson(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, maps.Values(users))
}

func GetUserByID(uuid string) (user, bool) {
	value, exists := users[uuid]
	return value, exists
}
