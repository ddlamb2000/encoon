// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package utils

import "github.com/google/uuid"

func GetNewUUID() string {
	return uuid.NewString()
}
