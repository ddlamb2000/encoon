// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package core

import "github.com/google/uuid"

func GetNewUUID() string {
	return uuid.NewString()
}
