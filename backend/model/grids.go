// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package model

import "fmt"

type Grid struct {
	Row
	Columns           []*Column       `json:"columns,omitempty"`
	Usages            []*Column       `json:"columnsUsage,omitempty"`
	Owners            map[string]bool `json:"-"`
	DefaultAccess     map[string]bool `json:"-"`
	ViewAccess        map[string]bool `json:"-"`
	EditAccess        map[string]bool `json:"-"`
	OwnerUuid         *string         `json:"-"`
	DefaultAccessUuid *string         `json:"-"`
	ViewAccessUuid    *string         `json:"-"`
	EditAccessUuid    *string         `json:"-"`
}

func GetNewGrid(uuid string) *Grid {
	grid := new(Grid)
	grid.Uuid = uuid
	grid.Owners = make(map[string]bool)
	grid.DefaultAccess = make(map[string]bool)
	grid.ViewAccess = make(map[string]bool)
	grid.EditAccess = make(map[string]bool)
	return grid
}

func (grid *Grid) SetDisplayString(dbName string) {
	if grid.Text1 != nil {
		grid.DisplayString = *grid.Text1
	}
}

func (grid *Grid) IsMetadata() bool {
	switch grid.Uuid {
	case UuidGrids:
		return true
	case UuidColumns:
		return true
	case UuidRelationships:
		return true
	case UuidMigrations:
		return true
	case UuidUsers:
		return true
	case UuidTransactions:
		return true
	default:
		return false
	}
	return false
}

func (grid *Grid) GetTableName() string {
	switch grid.Uuid {
	case UuidGrids:
		return "grids"
	case UuidColumns:
		return "columns"
	case UuidRelationships:
		return "relationships"
	case UuidMigrations:
		return "migrations"
	case UuidUsers:
		return "users"
	case UuidTransactions:
		return "transactions"
	default:
		return "rows"
	}
}

func (grid *Grid) CopyAccessToOtherGrid(otherGrid *Grid) {
	if grid.OwnerUuid != nil {
		otherGrid.Owners[*grid.OwnerUuid] = true
	}
	if grid.DefaultAccessUuid != nil {
		otherGrid.DefaultAccess[*grid.DefaultAccessUuid] = true
	}
	if grid.ViewAccessUuid != nil {
		otherGrid.ViewAccess[*grid.ViewAccessUuid] = true
	}
	if grid.EditAccessUuid != nil {
		otherGrid.EditAccess[*grid.EditAccessUuid] = true
	}
}

func (grid *Grid) GetViewEditAccessFlags(userUuid string) (canViewRows, canEditRows, canAddRows, canEditGrid bool) {
	switch {
	case grid.Owners[userUuid] && grid.Uuid == UuidTransactions:
		return true, false, false, false
	case grid.Owners[userUuid]:
		return true, true, true, true
	case grid.EditAccess[userUuid] || grid.DefaultAccess[UuidAccessLevelWriteAccess]:
		return true, true, true, false
	case grid.ViewAccess[userUuid] || grid.DefaultAccess[UuidAccessLevelReadAccess]:
		return true, false, false, false
	case grid.Uuid == UuidGrids || grid.Uuid == UuidRelationships || grid.Uuid == UuidColumns:
		return true, false, true, false
	case grid.Uuid == UuidAccessLevels || grid.Uuid == UuidUsers || grid.Uuid == UuidColumnTypes:
		return true, false, false, false
	}
	return false, false, false, false
}

func (grid *Grid) HasOwnership(userUuid string) bool {
	return grid.Owners[userUuid]
}

func (grid *Grid) GetRowsColumnDefinitions() string {
	columnDefinitions := ""
	for i := 1; i <= NumberOfTextFields; i++ {
		columnDefinitions += fmt.Sprintf(", text%d text", i)
	}
	for i := 1; i <= NumberOfIntFields; i++ {
		columnDefinitions += fmt.Sprintf(", int%d integer", i)
	}
	return columnDefinitions
}

func (grid *Grid) GetRowsQueryForExportDb() string {
	return grid.getRowsQueryColumnsForExportDb() + "FROM " + grid.GetTableName() + " ORDER BY created"
}

func (grid *Grid) getRowsQueryColumnsForExportDb() string {
	return "SELECT uuid, " +
		"gridUuid, " +
		"created, " +
		"createdBy, " +
		"updated, " +
		"updatedBy" +
		grid.getRowsColumnDefinitionsForExportDb() + ", " +
		"enabled, " +
		"revision "
}

func (grid *Grid) getRowsColumnDefinitionsForExportDb() string {
	columnDefinitions := ""
	switch grid.Uuid {
	case UuidGrids:
		columnDefinitions += ", text1, text2, text3"
	case UuidColumns:
		columnDefinitions += ", text1, text2, text3, int1"
	case UuidRelationships:
		columnDefinitions += ", text1, text2, text3, text4, text5"
	case UuidMigrations:
		columnDefinitions += ", text1, int1"
	case UuidUsers:
		columnDefinitions += ", text1, text2, text3, text4"
	case UuidTransactions:
		columnDefinitions += ", text1"
	default:
		for i := 1; i <= NumberOfTextFields; i++ {
			columnDefinitions += fmt.Sprintf(", text%d", i)
		}
		for i := 1; i <= NumberOfIntFields; i++ {
			columnDefinitions += fmt.Sprintf(", int%d", i)
		}
	}
	return columnDefinitions
}

func (grid *Grid) GetRowsQueryForSeedData() string {
	return "SELECT uuid, revision FROM " + grid.GetTableName() + " WHERE gridUuid = $1 AND uuid = $2"
}

func (grid *Grid) GetInsertStatementForSeedRowDb() string {
	return "INSERT INTO " + grid.GetTableName() +
		" (uuid, " +
		"revision, " +
		"created, " +
		"updated, " +
		"createdBy, " +
		"updatedBy, " +
		"enabled, " +
		"gridUuid" +
		grid.getRowsColumnDefinitionsForExportDb() + ") " +
		"VALUES ($1, " +
		"$2, " +
		"$3, " +
		"$4, " +
		"$5, " +
		"$6, " +
		"$7, " +
		"$8" +
		grid.getInsertStatementParametersForSeedRowDb() + ")"
}

func (grid *Grid) getInsertStatementParametersForSeedRowDb() string {
	parameters := ""
	switch grid.Uuid {
	case UuidGrids:
		parameters += ", $9, $10, $11"
	case UuidColumns:
		parameters += ", $9, $10, $11, $12"
	case UuidRelationships:
		parameters += ", $9, $10, $11, $12, $13"
	case UuidMigrations:
		parameters += ", $9, $10"
	case UuidUsers:
		parameters += ", $9, $10, $11, $12"
	case UuidTransactions:
		parameters += ", $9"
	default:
		parameterIndex := 9
		for i := 1; i <= NumberOfTextFields; i++ {
			parameters += fmt.Sprintf(", $%d", parameterIndex)
			parameterIndex += 1
		}
		for i := 1; i <= NumberOfIntFields; i++ {
			parameters += fmt.Sprintf(", $%d", parameterIndex)
			parameterIndex += 1
		}
	}
	return parameters
}

func (grid *Grid) GetInsertValuesForSeedRowDb(userUuid string, row *Row) []any {
	values := make([]any, 0)
	values = append(values, row.Uuid)
	values = append(values, row.Revision)
	values = append(values, row.Created)
	values = append(values, row.Updated)
	values = append(values, row.CreatedBy)
	values = append(values, row.UpdatedBy)
	values = append(values, row.Enabled)
	values = append(values, grid.Uuid)
	values = row.AppendRowValuesForSeedRowDb(values)
	return values
}

func (grid *Grid) GetUpdateValuesForSeedRowDb(userUuid string, row *Row) []any {
	values := make([]any, 0)
	values = append(values, grid.Uuid)
	values = append(values, row.Uuid)
	values = append(values, row.Revision)
	values = append(values, row.Updated)
	values = append(values, row.UpdatedBy)
	values = append(values, row.Enabled)
	values = row.AppendRowValuesForSeedRowDb(values)
	return values
}

func (grid *Grid) GetUpdateStatementForSeedRowDb() string {
	return "UPDATE " + grid.GetTableName() +
		" SET revision = $3, " +
		"updated = $4, " +
		"updatedBy = $5, " +
		"enabled = $6" +
		grid.getUpdateStatementParametersForSeedRowDb() +
		" WHERE gridUuid = $1 " +
		"AND uuid = $2"
}

func (grid *Grid) getUpdateStatementParametersForSeedRowDb() string {
	parameters := ""
	switch grid.Uuid {
	case UuidGrids:
		parameters += ", text1 = $7, text2 = $8, text3 = $9"
	case UuidColumns:
		parameters += ", text1 = $7, text2 = $8, text3 = $9, int1 = $10"
	case UuidUsers:
		parameters += ", text1 = $7, text2 = $8, text3 = $9, text4 = $10"
	case UuidRelationships:
		parameters += ", text1 = $7, text2 = $8, text3 = $9, text4 = $10, text5 = $11"
	default:
		parameterIndex := 7
		for i := 1; i <= NumberOfTextFields; i++ {
			parameters += fmt.Sprintf(", text%d = $%d", i, parameterIndex)
			parameterIndex += 1
		}
		for i := 1; i <= NumberOfIntFields; i++ {
			parameters += fmt.Sprintf(", int%d = $%d", i, parameterIndex)
			parameterIndex += 1
		}
	}
	return parameters
}
