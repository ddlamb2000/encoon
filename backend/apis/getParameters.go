// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2025

package apis

func getRequestParameters(dbName string, userUuid string, userName string, content requestContent) ApiParameters {
	return ApiParameters{
		DbName:   dbName,
		UserUuid: userUuid,
		UserName: userName,
		GridUuid: content.GridUuid,
		Uuid:     content.Uuid,
	}
}
