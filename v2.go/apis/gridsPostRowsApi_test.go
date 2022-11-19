// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"reflect"
	"testing"

	"d.lambert.fr/encoon/model"
	_ "github.com/lib/pq"
)

func TestGetInsertStatementForGridsApi(t *testing.T) {
	var grid model.Grid
	grid.Uuid = model.UuidUsers
	grid.Columns = append(grid.Columns, &model.Column{Name: "text2"})
	grid.Columns = append(grid.Columns, &model.Column{Name: "text5"})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int8"})
	got := getInsertStatementForGridsApi(&grid)
	expect := "INSERT INTO rows (uuid, version, created, updated, createdBy, updatedBy, enabled, gridUuid, text2, text5, int8)  VALUES ($1, 1, NOW(), NOW(), $2, $2, true, $3, $4, $5, $6)"
	if got != expect {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}

func TestGetInsertValuesForGridsApi(t *testing.T) {
	var grid model.Grid
	grid.Uuid = model.UuidUsers
	grid.Columns = append(grid.Columns, &model.Column{Name: "text2"})
	grid.Columns = append(grid.Columns, &model.Column{Name: "text5"})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int8"})
	uuid := "aaa"
	text2 := "yyy"
	text5 := "zzz"
	int8 := int64(10)
	row := model.Row{
		Uuid:  uuid,
		Text2: &text2,
		Text5: &text5,
		Int8:  &int8,
	}
	got := getInsertValuesForGridsApi("xxx", "yyy", &grid, &row)
	expect := []any{
		uuid,
		"xxx",
		"yyy",
		&text2,
		&text5,
		&int8,
	}
	if !reflect.DeepEqual(got, expect) {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}

func TestGetUpdateStatementForGridsApi(t *testing.T) {
	var grid model.Grid
	grid.Uuid = model.UuidUsers
	grid.Columns = append(grid.Columns, &model.Column{Name: "text2"})
	grid.Columns = append(grid.Columns, &model.Column{Name: "text5"})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int8"})
	got := getUpdateStatementForGridsApi(&grid)
	expect := "UPDATE rows SET version = version + 1, updated = NOW(), updatedBy = $3, text2 = $4, text5 = $5, int8 = $6 WHERE uuid = $1 and gridUuid = $2"
	if got != expect {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}

func TestGetUpdateValuesForGridsApi(t *testing.T) {
	var grid model.Grid
	grid.Uuid = model.UuidUsers
	grid.Columns = append(grid.Columns, &model.Column{Name: "text2"})
	grid.Columns = append(grid.Columns, &model.Column{Name: "text5"})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int8"})
	uuid := "aaa"
	text2 := "yyy"
	text5 := "zzz"
	int8 := int64(10)
	row := model.Row{
		Uuid:  uuid,
		Text2: &text2,
		Text5: &text5,
		Int8:  &int8,
	}
	got := getUpdateValuesForGridsApi("xxx", "yyy", &grid, &row)
	expect := []any{
		uuid,
		"yyy",
		"xxx",
		&text2,
		&text5,
		&int8,
	}
	if !reflect.DeepEqual(got, expect) {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}
