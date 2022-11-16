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
	got := getInsertStatementForGridsApi()
	expect := "INSERT INTO rows (uuid, version, created, updated, createdBy, updatedBy, enabled, gridUuid, text1, text2, text3, text4, int1, int2, int3, int4)  VALUES ($1, 1, NOW(), NOW(), $2, $2, true, $3, $4, $5, $6, $7, $8, $9, $10, $11)"
	if got != expect {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}

func TestGetInsertValuesForGridsApi(t *testing.T) {
	uuid := "aaa"
	text1 := "zzz"
	int1 := int64(10)
	row := model.Row{
		Uuid:  uuid,
		Text1: &text1,
		Text2: &text1,
		Text3: &text1,
		Text4: &text1,
		Int1:  &int1,
		Int2:  &int1,
		Int3:  &int1,
		Int4:  &int1,
	}
	got := getInsertValuesForGridsApi("xxx", "yyy", row)
	expect := []any{
		uuid,
		"xxx",
		"yyy",
		&text1,
		&text1,
		&text1,
		&text1,
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
	got := getUpdateStatementForGridsApi()
	expect := "UPDATE rows SET version = version + 1, updated = NOW(), updatedBy = $3, text1 = $4, text2 = $5, text3 = $6, text4 = $7, int1 = $8, int2 = $9, int3 = $10, int4 = $11 WHERE uuid = $1 and gridUuid = $2"
	if got != expect {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}

func TestGetUpdateValuesForGridsApi(t *testing.T) {
	uuid := "aaa"
	text1 := "zzz"
	int1 := int64(10)
	row := model.Row{
		Uuid:  uuid,
		Text1: &text1,
		Text2: &text1,
		Text3: &text1,
		Text4: &text1,
		Int1:  &int1,
		Int2:  &int1,
		Int3:  &int1,
		Int4:  &int1,
	}
	got := getUpdateValuesForGridsApi("xxx", "yyy", row)
	expect := []any{
		uuid,
		"yyy",
		"xxx",
		&text1,
		&text1,
		&text1,
		&text1,
		&int1,
		&int1,
		&int1,
		&int1,
	}
	if !reflect.DeepEqual(got, expect) {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}
