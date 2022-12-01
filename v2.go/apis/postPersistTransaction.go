// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"d.lambert.fr/encoon/model"
	"d.lambert.fr/encoon/utils"
)

func postInsertTransaction(r apiRequestParameters) error {
	query := getInsertStatementForTransaction()
	parms := getInsertValuesForTransaction(r.userUuid, r.transaction)
	r.trace("postInsertTransaction() - query=%s, parms=%s", query, parms)
	if err := r.execContext(query, parms...); err != nil {
		_ = r.rollbackTransaction()
		return r.logAndReturnError("Insert transaction error: %v.", err)
	}
	r.log("Transaction [%s] inserted.", r.transaction.Uuid)
	return nil
}

// function is available for mocking
var getInsertStatementForTransaction = func() string {
	insertStr := "INSERT INTO transactions " +
		"(uuid, " +
		"revision, " +
		"created, " +
		"updated, " +
		"createdBy, " +
		"updatedBy, " +
		"enabled, " +
		"gridUuid, " +
		"text1" +
		") "
	valueStr := " VALUES ($1, " +
		"1, " +
		"NOW(), " +
		"NOW(), " +
		"$2, " +
		"$2, " +
		"true, " +
		"$3, " +
		"$4" +
		")"
	return insertStr + valueStr
}

func getInsertValuesForTransaction(userUuid string, row *model.Row) []any {
	values := make([]any, 0)
	values = append(values, row.Uuid)
	values = append(values, userUuid)
	values = append(values, model.UuidTransactions)
	values = append(values, row.Text1)
	return values
}

func postInsertTransactionReferenceRow(r apiRequestParameters, grid *model.Grid, row *model.Row, relationship string) error {
	query := getInsertStatementForTransactionReferenceRow()
	parms := getInsertStatementParametersForTransactionReferenceRow(r, grid, row, relationship)
	r.trace("postInsertTransactionReferenceRow(%s) - query=%s ; parms=%s", row, query, parms)
	if err := r.execContext(query, parms...); err != nil {
		return r.logAndReturnError("Insert transaction referenced row error: %v.", err)
	}
	r.log("Referenced row [%v] inserted into transaction %s.", row, r.transaction.Uuid)
	return nil
}

// function is available for mocking
var getInsertStatementForTransactionReferenceRow = func() string {
	return "INSERT INTO relationships (uuid, " +
		"revision, " +
		"created, " +
		"updated, " +
		"createdBy, " +
		"updatedBy, " +
		"enabled, " +
		"gridUuid, " +
		"text1, " +
		"text2, " +
		"text3, " +
		"text4, " +
		"text5) " +
		"VALUES ($1, " +
		"1, " +
		"NOW(), " +
		"NOW(), " +
		"$2, " +
		"$2, " +
		"true, " +
		"$3, " +
		"$4, " +
		"$5, " +
		"$6, " +
		"$7, " +
		"$8)"
}

func getInsertStatementParametersForTransactionReferenceRow(r apiRequestParameters, grid *model.Grid, row *model.Row, relationship string) []any {
	parameters := make([]any, 0)
	parameters = append(parameters, utils.GetNewUUID())
	parameters = append(parameters, r.userUuid)
	parameters = append(parameters, model.UuidRelationships)
	parameters = append(parameters, relationship)
	parameters = append(parameters, model.UuidTransactions)
	parameters = append(parameters, r.transaction.Uuid)
	parameters = append(parameters, grid.Uuid)
	parameters = append(parameters, row.Uuid)
	return parameters
}
