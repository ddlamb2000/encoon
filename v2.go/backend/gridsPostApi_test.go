// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import (
	"reflect"
	"testing"

	_ "github.com/lib/pq"
)

func TestGetInsertStatementForGridsApi(t *testing.T) {
	got := getInsertStatementForGridsApi()
	expected := "INSERT INTO rows (uuid, version, created, updated, createdBy, updatedBy, enabled, gridUuid, text01, text02, text03, text04, int01, int02, int03, int04)  VALUES ($1, 1, NOW(), NOW(), $2, $2, true, $3, $4, $5, $6, $7, $8, $9, $10, $11)"
	if got != expected {
		t.Errorf(`Got %v instead of %v.`, got, expected)
	}
}

func TestGetInsertValuesForGridsApi(t *testing.T) {
	uuid := "aaa"
	text01 := "zzz"
	int01 := int64(10)
	row := Row{
		Uuid:   uuid,
		Text01: &text01,
		Text02: &text01,
		Text03: &text01,
		Text04: &text01,
		Int01:  &int01,
		Int02:  &int01,
		Int03:  &int01,
		Int04:  &int01,
	}
	got := getInsertValuesForGridsApi("xxx", "yyy", row)
	expected := []any{
		uuid,
		"xxx",
		"yyy",
		&text01,
		&text01,
		&text01,
		&text01,
		&int01,
		&int01,
		&int01,
		&int01,
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf(`Got %v instead of %v.`, got, expected)
	}
}

func TestGetUpdateStatementForGridsApi(t *testing.T) {
	got := getUpdateStatementForGridsApi()
	expected := "UPDATE rows SET version = version + 1, updated = NOW(), updatedBy = $3, text01 = $4, text02 = $5, text03 = $6, text04 = $7, int01 = $8, int02 = $9, int03 = $10, int04 = $11 WHERE uuid = $1 and gridUuid = $2"
	if got != expected {
		t.Errorf(`Got %v instead of %v.`, got, expected)
	}
}

func TestGetUpdateValuesForGridsApi(t *testing.T) {
	uuid := "aaa"
	text01 := "zzz"
	int01 := int64(10)
	row := Row{
		Uuid:   uuid,
		Text01: &text01,
		Text02: &text01,
		Text03: &text01,
		Text04: &text01,
		Int01:  &int01,
		Int02:  &int01,
		Int03:  &int01,
		Int04:  &int01,
	}
	got := getUpdateValuesForGridsApi("xxx", "yyy", row)
	expected := []any{
		uuid,
		"yyy",
		"xxx",
		&text01,
		&text01,
		&text01,
		&text01,
		&int01,
		&int01,
		&int01,
		&int01,
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf(`Got %v instead of %v.`, got, expected)
	}
}
