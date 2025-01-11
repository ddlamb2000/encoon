// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2025

package kafka

import (
	"d.lambert.fr/encoon/apis"
)

func getParameters(dbName string, userUuid string, userName string, content requestContent) apis.ApiParameters {
	return apis.ApiParameters{
		DbName:   dbName,
		UserUuid: userUuid,
		UserName: userName,
		GridUuid: content.GridUuid,
		Uuid:     content.Uuid,
	}
}
