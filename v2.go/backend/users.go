// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import (
	"fmt"

	"d.lambert.fr/encoon/utils"
)

type user struct {
	entity

	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

var (
	users = make(map[string]user)
)

func (u user) String() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

func loadUsers() {
	utils.Log("Loading users.")

	(&user{
		entity:    entity{Uuid: "c788a76d-4aa6-4073-8904-35a9b99a3289", Version: 1, Enabled: true},
		Email:     "root@encoon.com",
		FirstName: "Root",
		LastName:  "Encoon"}).add()

	(&user{
		entity:    entity{Uuid: "bced42a2-6ddd-4023-ad40-0d46962b7872", Version: 1, Enabled: true},
		Email:     "system@encoon.com",
		FirstName: "System",
		LastName:  "Encoon"}).add()

	(&user{
		entity: entity{Uuid: "67b560b9-63ff-4fed-9b64-26c7f86e540c"},
		Email:  "none@encoon.com"}).add()

	utils.Log("Users loaded.")
}

func (v *user) add() {
	users[v.entity.Uuid] = *v
}

func GetUserByID(uuid string) (user, bool) {
	value, exists := users[uuid]
	return value, exists
}
