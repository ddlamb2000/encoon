// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import "d.lambert.fr/encoon/model"

func getAuditsForRow(r apiRequest, grid *model.Grid, uuid string) ([]*model.Audit, error) {
	r.trace("getAuditsForRow(%s, %s)", grid, uuid)
	t := r.startTiming()
	defer r.stopTiming("getAuditsForRow()", t)
	query := getAuditsQueryForRow(grid, uuid)
	parms := getAuditsQueryParametersForRow(grid.Uuid, uuid)
	r.trace("getAuditsForRow(%s, %s) - query=%s ; parms=%s", grid, uuid, query, parms)
	set, err := r.db.QueryContext(r.ctx, query, parms...)
	if err != nil {
		return nil, r.logAndReturnError("Error when querying audits: %v.", err)
	}
	defer set.Close()
	audits := make([]*model.Audit, 0)
	for set.Next() {
		audit := model.GetNewAudit()
		if err := set.Scan(getAuditsQueryOutputForRow(audit)...); err != nil {
			return nil, r.logAndReturnError("Error when scanning audits: %v.", err)
		}
		audit.SetActionName()
		audits = append(audits, audit)
	}
	return audits, nil
}

// function is available for mocking
var getAuditsQueryForRow = func(grid *model.Grid, uuid string) string {
	return "SELECT transactions.uuid, " +
		"transactions.created, " +
		"transactions.createdBy, " +
		"createdBy.text1, " +
		"relCreated.text1 " +
		"FROM transactions " +

		"INNER JOIN relationships relCreated " +
		"ON relCreated.gridUuid = $3 " +
		"AND relCreated.text2 = transactions.gridUuid " +
		"AND relCreated.text3 = transactions.uuid " +
		"AND relCreated.text4 = $4 " +
		"AND relCreated.text5 = $5 " +
		"AND relCreated.enabled = true " +

		"LEFT OUTER JOIN users createdBy " +
		"ON createdBy.gridUuid = $1 " +
		"AND createdBy.uuid = transactions.createdBy " +
		"AND createdBy.enabled = true " +

		"WHERE transactions.griduuid = $2 " +
		"ORDER BY transactions.created DESC"
}

var getAuditsQueryParametersForRow = func(gridUuid, uuid string) []any {
	parameters := make([]any, 0)
	parameters = append(parameters, model.UuidUsers)
	parameters = append(parameters, model.UuidTransactions)
	parameters = append(parameters, model.UuidRelationships)
	parameters = append(parameters, gridUuid)
	parameters = append(parameters, uuid)
	return parameters
}

// function is available for mocking
var getAuditsQueryOutputForRow = func(audit *model.Audit) []any {
	output := make([]any, 0)
	output = append(output, &audit.Uuid)
	output = append(output, &audit.Created)
	output = append(output, &audit.CreatedBy)
	output = append(output, &audit.CreatedByName)
	output = append(output, &audit.ColumnName)
	return output
}
