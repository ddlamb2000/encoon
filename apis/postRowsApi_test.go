// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package apis

import (
	"reflect"
	"testing"

	"d.lambert.fr/encoon/model"
	_ "github.com/lib/pq"
)

func TestGetInsertStatementForGridsApi(t *testing.T) {
	grid := model.GetNewGrid()
	grid.Uuid = model.UuidUsers
	grid.Columns = append(grid.Columns, &model.Column{Name: "text2", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "text5", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int8", Owned: true})
	got := getInsertStatementForGridsApi(grid)
	expect := "INSERT INTO users (uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, text2, text5, int8) VALUES ($1, 1, NOW(), NOW(), $2, $2, true, $3, $4, $5, $6)"
	if got != expect {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}

func TestGetInsertValuesForGridsApi(t *testing.T) {
	grid := model.GetNewGrid()
	grid.Uuid = model.UuidUsers
	grid.Columns = append(grid.Columns, &model.Column{Name: "text2", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "text5", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int8", Owned: true})
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
	got := getInsertValuesForGridsApi("xxx", grid, &row)
	expect := []any{
		uuid,
		"xxx",
		model.UuidUsers,
		&text2,
		&text5,
		&int8,
	}
	if !reflect.DeepEqual(got, expect) {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}

func TestGetInsertValuesForGridsApi2(t *testing.T) {
	grid := model.GetNewGrid()
	grid.Uuid = model.UuidUsers
	grid.Columns = append(grid.Columns, &model.Column{Name: "text1", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "text2", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "text3", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "text4", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "text6", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "text7", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "text8", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "text9", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "text10", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int1", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int2", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int3", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int4", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int5", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int6", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int7", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int9", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int10", Owned: true})
	uuid := "aaa"
	text1 := "yyy"
	int1 := int64(10)
	row := model.Row{
		Uuid:   uuid,
		Text1:  &text1,
		Text2:  &text1,
		Text3:  &text1,
		Text4:  &text1,
		Text6:  &text1,
		Text7:  &text1,
		Text8:  &text1,
		Text9:  &text1,
		Text10: &text1,
		Int1:   &int1,
		Int2:   &int1,
		Int3:   &int1,
		Int4:   &int1,
		Int5:   &int1,
		Int6:   &int1,
		Int7:   &int1,
		Int9:   &int1,
		Int10:  &int1,
	}
	got := getInsertValuesForGridsApi("xxx", grid, &row)
	expect := []any{
		uuid,
		"xxx",
		model.UuidUsers,
		&text1,
		&text1,
		&text1,
		&text1,
		&text1,
		&text1,
		&text1,
		&text1,
		&text1,
		&int1,
		&int1,
		&int1,
		&int1,
		&int1,
		&int1,
		&int1,
		&int1,
		&int1,
	}
	if !reflect.DeepEqual(got, expect) {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}

func TestGetUpdateStatementForGridsApi(t *testing.T) {
	grid := model.GetNewGrid()
	grid.Uuid = model.UuidUsers
	grid.Columns = append(grid.Columns, &model.Column{Name: "text2", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "text5", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int8", Owned: true})
	got := getUpdateStatementForGridsApi(grid)
	expect := "UPDATE users SET revision = revision + 1, updated = NOW(), updatedBy = $3, text2 = $4, text5 = $5, int8 = $6 WHERE gridUuid = $1 AND uuid = $2"
	if got != expect {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}

func TestGetUpdateValuesForGridsApi(t *testing.T) {
	grid := model.GetNewGrid()
	grid.Uuid = model.UuidUsers
	grid.Columns = append(grid.Columns, &model.Column{Name: "text2", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "text5", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int8", Owned: true})
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
	got := getUpdateValuesForGridsApi("xxx", grid, &row)
	expect := []any{
		model.UuidUsers,
		uuid,
		"xxx",
		&text2,
		&text5,
		&int8,
	}
	if !reflect.DeepEqual(got, expect) {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}
